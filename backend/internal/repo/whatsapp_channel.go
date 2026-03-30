package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WhatsAppChannelRecord struct {
	ID             string
	AccountID      string
	PhoneNumber    string
	JID            string
	Provider       string
	ProviderConfig []byte
	QRCode         string
	Connected      bool
	ConnectedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WhatsAppChannelRepository struct {
	db *pgxpool.Pool
}

func NewWhatsAppChannelRepository(db *pgxpool.Pool) *WhatsAppChannelRepository {
	return &WhatsAppChannelRepository{db: db}
}

func (r *WhatsAppChannelRepository) Create(ctx context.Context, ch *WhatsAppChannelRecord) error {
	q := `INSERT INTO "wzChannelsWhatsapp" ("id","accountId","phoneNumber","jid","provider","providerConfig","qrCode","connected","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(ctx, q,
		ch.ID, ch.AccountID, ch.PhoneNumber, ch.JID, ch.Provider,
		ch.ProviderConfig, ch.QRCode, ch.Connected, ch.CreatedAt, ch.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert whatsapp channel: %w", err)
	}
	return nil
}

func (r *WhatsAppChannelRepository) FindByID(ctx context.Context, id string) (*WhatsAppChannelRecord, error) {
	q := `SELECT "id","accountId",COALESCE("phoneNumber",''),COALESCE("jid",''),"provider","providerConfig",
		  COALESCE("qrCode",''),"connected","connectedAt","createdAt","updatedAt"
		  FROM "wzChannelsWhatsapp" WHERE "id"=$1`
	var ch WhatsAppChannelRecord
	err := r.db.QueryRow(ctx, q, id).Scan(
		&ch.ID, &ch.AccountID, &ch.PhoneNumber, &ch.JID, &ch.Provider,
		&ch.ProviderConfig, &ch.QRCode, &ch.Connected, &ch.ConnectedAt,
		&ch.CreatedAt, &ch.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("whatsapp channel not found: %w", err)
	}
	return &ch, nil
}

func (r *WhatsAppChannelRepository) UpdateQRCode(ctx context.Context, id string, qr string) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzChannelsWhatsapp" SET "qrCode"=$1 WHERE "id"=$2`, qr, id)
	return err
}

func (r *WhatsAppChannelRepository) UpdateJID(ctx context.Context, id string, jid string) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzChannelsWhatsapp" SET "jid"=$1, "connected"=true, "connectedAt"=NOW() WHERE "id"=$2`, jid, id)
	return err
}

func (r *WhatsAppChannelRepository) SetConnected(ctx context.Context, id string, connected bool) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzChannelsWhatsapp" SET "connected"=$1 WHERE "id"=$2`, connected, id)
	return err
}

func (r *WhatsAppChannelRepository) ClearDevice(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzChannelsWhatsapp" SET "connected"=false,"jid"='','qrCode'='' WHERE "id"=$1`, id)
	return err
}

func (r *WhatsAppChannelRepository) FindJIDByID(ctx context.Context, id string) (string, error) {
	var jid string
	err := r.db.QueryRow(ctx, `SELECT COALESCE("jid",'') FROM "wzChannelsWhatsapp" WHERE "id"=$1`, id).Scan(&jid)
	if err != nil {
		return "", fmt.Errorf("whatsapp channel not found: %w", err)
	}
	return jid, nil
}

func (r *WhatsAppChannelRepository) FindByJID(ctx context.Context, jid string) (*WhatsAppChannelRecord, error) {
	q := `SELECT "id","accountId",COALESCE("phoneNumber",''),COALESCE("jid",''),"provider","providerConfig",
		  COALESCE("qrCode",''),"connected","connectedAt","createdAt","updatedAt"
		  FROM "wzChannelsWhatsapp" WHERE "jid"=$1`
	var ch WhatsAppChannelRecord
	err := r.db.QueryRow(ctx, q, jid).Scan(
		&ch.ID, &ch.AccountID, &ch.PhoneNumber, &ch.JID, &ch.Provider,
		&ch.ProviderConfig, &ch.QRCode, &ch.Connected, &ch.ConnectedAt,
		&ch.CreatedAt, &ch.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("whatsapp channel not found for jid %s: %w", jid, err)
	}
	return &ch, nil
}

func (r *WhatsAppChannelRepository) FindQRByID(ctx context.Context, id string) (string, error) {
	var qr string
	err := r.db.QueryRow(ctx, `SELECT COALESCE("qrCode",'') FROM "wzChannelsWhatsapp" WHERE "id"=$1`, id).Scan(&qr)
	if err != nil {
		return "", fmt.Errorf("whatsapp channel not found: %w", err)
	}
	return qr, nil
}

func (r *WhatsAppChannelRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzChannelsWhatsapp" WHERE "id"=$1`, id)
	return err
}
