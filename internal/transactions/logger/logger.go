package logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"orange-go/internal/storage"
	"orange-go/internal/transactions"
)

type FileTransactionLogger struct {
	events        chan<- transactions.Event // Канал только для записи; для передачи событий
	errors        <-chan error              // Канал только для чтения; для приема ошибок
	lastSequence  uint64                    // Последний использованный порядковый номер
	file          *os.File                  // Местоположение файла журнала
	memoryStorage storage.IMemoryStorage
	done          chan struct{} // сигнал завершения писателя
	closed        bool
	mu            sync.Mutex
}

func NewFileTransactionLogger(filename string, memoryStorage storage.IMemoryStorage) (transactions.ITransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{
		file:          file,
		memoryStorage: memoryStorage,
		done:          make(chan struct{}),
	}, nil
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- transactions.Event{EventType: transactions.EventPut, Key: key, Value: value}
}
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- transactions.Event{EventType: transactions.EventDelete, Key: key}
}
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) GetMemoryStorage() storage.IMemoryStorage {
	return l.memoryStorage
}

func (l *FileTransactionLogger) Run() {
	events := make(chan transactions.Event, 16) // Создать канал событий
	l.events = events
	errors := make(chan error, 1) // Создать канал ошибок
	l.errors = errors
	go func() {
		defer close(l.done)     // <- сообщаем Close(), что писатель завершился
		for e := range events { // Извлечь следующее событие Event
			l.lastSequence++       // Увеличить порядковый номер
			_, err := fmt.Fprintf( // Записать событие в журнал
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan transactions.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)       // Создать Scanner для чтения l.file
	outEvent := make(chan transactions.Event) // Небуферизованный канал событий
	outError := make(chan error, 1)           // Буферизованный канал ошибок
	go func() {
		var e transactions.Event
		defer close(outEvent) // Закрыть каналы
		defer close(outError) // по завершении сопрограммы
		for scanner.Scan() {
			line := scanner.Text()
			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}
			// Проверка целостности!
			// Порядковые номера последовательно увеличиваются?
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			l.lastSequence = e.Sequence // Запомнить последний использованный
			// порядковый номер
			outEvent <- e // Отправить событие along
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent, outError
}

func InitializeTransactionLog(tl transactions.ITransactionLogger) error {
	events, errors := tl.ReadEvents()

	var (
		e   transactions.Event
		ok  = true
		err error
	)

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transactions.EventDelete:
				err = tl.GetMemoryStorage().Delete(e.Key)
			case transactions.EventPut:
				err = tl.GetMemoryStorage().Put(e.Key, e.Value)
			}
		}
	}
	tl.Run()
	return err
}

// Close корректно завершает горутину и закрывает файл
func (l *FileTransactionLogger) Close() error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return nil
	}
	l.closed = true
	// закрываем канал событий, чтобы писатель завершился
	if l.events != nil {
		close(l.events)
	}
	l.mu.Unlock()

	// ждём завершения писателя
	<-l.done

	// убедимся, что всё записано и закроем файл
	if err := l.file.Sync(); err != nil {
		_ = l.file.Close()
		return fmt.Errorf("sync log file: %w", err)
	}
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("close log file: %w", err)
	}
	return nil
}
