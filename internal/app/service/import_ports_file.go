package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/FurmanovD/portsvc/internal/pkg/models"
)

const (
	maxPortsToAddInBatch = 20
)

// ImportPortsFile reads the file with ports information and upserts them into the DB in a batch.
func (s *serviceImpl) ImportPortsFile(ctx context.Context, portsFilePath string) error {

	portsCh := make(chan *models.Port)
	doneCh := make(chan struct{})
	// start ports reader.
	go s.getPortsFromFile(portsFilePath, portsCh, doneCh)

	// use buffer to save ports using batch insert.
	portsBuff := make([]*models.Port, maxPortsToAddInBatch)
	currentIdx := -1
	flushBuffer := false
	continueProcessing := true
	for continueProcessing {
		flushBuffer = false

		select {
		case port := <-portsCh:
			currentIdx++
			portsBuff[currentIdx] = port

			flushBuffer = currentIdx >= maxPortsToAddInBatch-1
		case <-doneCh:
			flushBuffer = true
			continueProcessing = false
		}

		if flushBuffer {
			//
			if err := s.savePortsInfo(ctx, portsBuff[:currentIdx+1]); err != nil {
				// TODO implement persistent buffer(queue?) to retry info saving.
				s.log.Errorf("saving ports %+v into DB error: %v", portsBuff, err)
				return err
			}
			// and cleanup the buffer.
			currentIdx = -1
			for i := range portsBuff {
				portsBuff[i] = nil
			}
		}
	}

	return nil
}

func (s *serviceImpl) getPortsFromFile(portsFilePath string, out chan *models.Port, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	f, err := os.Open(portsFilePath)
	if err != nil {
		s.log.Errorf("failed to open incoming file '%s': %v", portsFilePath, err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	// read open bracket
	t, err := dec.Token()
	fmt.Printf("\n\n%+v\n\n", t)
	if err != nil {
		s.log.Errorf("failed to decode ports JSON opening bracket: %v", err)
		return
	}

	// while the file contains port info.
	for dec.More() {
		portID, err := dec.Token()
		if err != nil {
			s.log.Errorf("failed to decode portID element: %v", err)
			return
		}

		// decode the object.
		var p models.Port
		err = dec.Decode(&p)
		if err != nil {
			s.log.Errorf("failed to decode port object: %v", err)
			return
		}

		// and set correct portID field.
		p.PortID = fmt.Sprintf("%s", portID)

		// send the constructed object to an out channel.
		out <- &p
	}

	// read closing bracket.
	_, err = dec.Token()
	if err != nil {
		s.log.Errorf("failed to decode closing JSON bracket: %v", err)
		return
	}
}

func (s *serviceImpl) savePortsInfo(
	ctx context.Context,
	ports []*models.Port,

) error {
	if len(ports) == 0 {
		return nil
	}

	// TODO add a transaction parameter to all repository methods, create it here
	// and pass to every methods of saving info to all other tables:
	// `ports_unlocks`, `ports_aliases`, (?)'regions'

	if err := s.db.Ports.SetPortsInfo(ctx, s.converter.ModelToDBPorts(ports)); err != nil {
		// TODO implement persistent buffer(queue?) to retry info saving.
		return err
	}

	// TODO add all the rest info in transaction.
	// var aliases []string
	// var unlocs []string
	// var regions [] type ?
	// for _, port := range ports {
	// 	aliases = append(aliases, port.Alias...)
	// 	unlocs = append(unlocs, port.Unlocs...)
	// }

	// TODO save aliases, unlocks, regions.

	// TODO commit transaction.

	return nil
}
