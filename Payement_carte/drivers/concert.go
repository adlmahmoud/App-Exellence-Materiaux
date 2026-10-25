package drivers

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"exellence-materiaux/Payement_carte/models"
)

const (
	ByteSTX = 0x02
	ByteETX = 0x03
	ByteEOT = 0x04
	ByteENQ = 0x05
	ByteACK = 0x06
	ByteNAK = 0x15
	ByteFS  = 0x1C
)

func calculerLRC(data []byte) byte {
	var lrc byte
	for _, b := range data {
		lrc ^= b
	}
	return lrc
}

func construireTrameConcert(caisse int, montantCentimes int64) []byte {
	caisseStr := fmt.Sprintf("%02d", caisse%100)
	montantStr := fmt.Sprintf("%08d", montantCentimes)
	typeTrans := "1"
	devise := "978"

	corps := []byte(caisseStr + montantStr + typeTrans + devise)
	lrc := calculerLRC(append(corps, ByteETX))

	trame := make([]byte, 0, len(corps)+3)
	trame = append(trame, ByteSTX)
	trame = append(trame, corps...)
	trame = append(trame, ByteETX, lrc)
	return trame
}

func dialoguerConcert(ctx context.Context, rw io.ReadWriter, caisse int, montantCentimes int64, timeout time.Duration) (models.ReponsePaiement, error) {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	dateStr := time.Now().Format("2006-01-02 15:04:05")
	montantEuros := float64(montantCentimes) / 100.0

	errCh := make(chan error, 1)
	respCh := make(chan models.ReponsePaiement, 1)

	go func() {
		buf := make([]byte, 512)

		if _, err := rw.Write([]byte{ByteENQ}); err != nil {
			errCh <- fmt.Errorf("envoi ENQ: %w", err)
			return
		}

		n, err := rw.Read(buf)
		if err != nil || n == 0 || buf[0] != ByteACK {
			errCh <- fmt.Errorf("attente ACK après ENQ invalide ou timeout")
			return
		}

		trameDemande := construireTrameConcert(caisse, montantCentimes)
		if _, err := rw.Write(trameDemande); err != nil {
			errCh <- fmt.Errorf("envoi trame paiement: %w", err)
			return
		}

		n, err = rw.Read(buf)
		if err != nil || n == 0 || buf[0] != ByteACK {
			errCh <- fmt.Errorf("attente ACK après trame invalide")
			return
		}

		if _, err := rw.Write([]byte{ByteEOT}); err != nil {
			errCh <- fmt.Errorf("envoi EOT: %w", err)
			return
		}

		var b byte
		var attenduENQ bool
		oneByte := make([]byte, 1)

		for !attenduENQ {
			if _, err := rw.Read(oneByte); err != nil {
				errCh <- fmt.Errorf("attente réponse TPE: %w", err)
				return
			}
			b = oneByte[0]
			if b == ByteENQ {
				attenduENQ = true
			}
		}

		if _, err := rw.Write([]byte{ByteACK}); err != nil {
			errCh <- fmt.Errorf("acquittement réponse TPE: %w", err)
			return
		}

		var reponseRaw []byte
		var dansMessage bool

		for {
			if _, err := rw.Read(oneByte); err != nil {
				errCh <- fmt.Errorf("lecture données TPE: %w", err)
				return
			}
			b = oneByte[0]
			if b == ByteSTX {
				dansMessage = true
				reponseRaw = reponseRaw[:0]
				continue
			}
			if dansMessage {
				if b == ByteETX {
					if _, err := rw.Read(oneByte); err == nil {
					}
					break
				}
				reponseRaw = append(reponseRaw, b)
			}
		}

		_, _ = rw.Write([]byte{ByteACK})
		_, _ = rw.Read(oneByte)

		rep := interpreterReponseConcert(reponseRaw, montantCentimes, montantEuros, dateStr)
		respCh <- rep
	}()

	select {
	case <-ctx.Done():
		return models.ReponsePaiement{
			Statut:          models.StatutAnnule,
			CodeReponse:     "99",
			Message:         "Paiement annulé depuis la caisse",
			MontantCentimes: montantCentimes,
			MontantEuros:    montantEuros,
			DateHeure:       dateStr,
		}, ctx.Err()
	case err := <-errCh:
		return models.ReponsePaiement{
			Statut:          models.StatutErreur,
			CodeReponse:     "ER",
			Message:         "Erreur de communication avec le TPE : " + err.Error(),
			MontantCentimes: montantCentimes,
			MontantEuros:    montantEuros,
			DateHeure:       dateStr,
		}, err
	case resp := <-respCh:
		return resp, nil
	case <-time.After(timeout):
		return models.ReponsePaiement{
			Statut:          models.StatutAnnule,
			CodeReponse:     "TO",
			Message:         "Délai d'attente du TPE dépassé (Timeout)",
			MontantCentimes: montantCentimes,
			MontantEuros:    montantEuros,
			DateHeure:       dateStr,
		}, fmt.Errorf("timeout TPE")
	}
}

func interpreterReponseConcert(corps []byte, montantCentimes int64, montantEuros float64, dateStr string) models.ReponsePaiement {
	texte := string(corps)

	codeReponse := "00"
	if len(texte) >= 12 {
		codeReponse = strings.TrimSpace(texte[10:12])
	}

	numAuth := ""
	if len(texte) >= 18 {
		numAuth = strings.TrimSpace(texte[12:18])
	}
	if numAuth == "" {
		numAuth = "AUT-OK"
	}

	switch codeReponse {
	case "00":
		return models.ReponsePaiement{
			Statut:             models.StatutValide,
			CodeReponse:        "00",
			Message:            "Paiement accepté par le centre bancaire",
			MontantCentimes:    montantCentimes,
			MontantEuros:       montantEuros,
			NumeroAutorisation: numAuth,
			NumeroTransaction:  fmt.Sprintf("TRX-%d", time.Now().UnixNano()%100000000),
			TypeCarte:          "CB / Sans contact",
			DateHeure:          dateStr,
		}
	case "05", "51", "04":
		return models.ReponsePaiement{
			Statut:          models.StatutRefuse,
			CodeReponse:     codeReponse,
			Message:         "Paiement refusé par la banque (Solde insuffisant ou plafond dépassé)",
			MontantCentimes: montantCentimes,
			MontantEuros:    montantEuros,
			DateHeure:       dateStr,
		}
	default:
		return models.ReponsePaiement{
			Statut:          models.StatutAnnule,
			CodeReponse:     codeReponse,
			Message:         "Transaction annulée sur le terminal (Touche rouge ou abandon)",
			MontantCentimes: montantCentimes,
			MontantEuros:    montantEuros,
			DateHeure:       dateStr,
		}
	}
}
