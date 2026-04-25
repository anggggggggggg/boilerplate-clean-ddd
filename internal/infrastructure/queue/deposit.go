package queue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

// TaskDepositCreated adalah nama task yang di-enqueue saat deposit berhasil dibuat.
// Nama ini digunakan oleh producer (EnqueueDepositCreated) dan consumer (ProcessDepositCreated).
const TaskDepositCreated = "deposit:created"

// depositCreatedPayload adalah struktur payload task deposit:created.
type depositCreatedPayload struct {
	DepositID string `json:"deposit_id"`
}

// depositQueue mengimplementasikan usecase.DepositQueue menggunakan Asynq.
type depositQueue struct {
	client *asynq.Client
}

// NewDepositQueue membuat instance depositQueue baru.
func NewDepositQueue(client *asynq.Client) *depositQueue {
	return &depositQueue{client: client}
}

// EnqueueDepositCreated memproduksi task ke antrian saat deposit baru dibuat.
// Task ini kemudian diproses oleh worker secara asynchronous.
func (q *depositQueue) EnqueueDepositCreated(ctx context.Context, depositID string) error {
	payload, err := json.Marshal(depositCreatedPayload{DepositID: depositID})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskDepositCreated, payload)

	// EnqueueContext mendukung context untuk cancellation
	info, err := q.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Printf("[Queue] task %s di-enqueue: id=%s, queue=%s", TaskDepositCreated, info.ID, info.Queue)
	return nil
}

// ---------------------------------------------------------------------------
// Worker — Task Processor
// ---------------------------------------------------------------------------

// ProcessDepositCreated adalah handler yang dipanggil oleh Asynq worker
// setiap kali ada task "deposit:created" di antrian.
// Di sinilah side-effect dijalankan: kirim email, push notifikasi, update statistik, dll.
func ProcessDepositCreated(ctx context.Context, t *asynq.Task) error {
	var payload depositCreatedPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Simulasi proses: di aplikasi nyata bisa kirim notifikasi, update ledger, dll.
	log.Printf("[Worker] memproses deposit:created | deposit_id=%s", payload.DepositID)

	return nil
}
