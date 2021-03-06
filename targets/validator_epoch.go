package targets

import (
	"context"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"

	"github.com/Chainflow/oasis-mission-control/config"
)

// GetValEpoch returns the work epoch number
func GetValEpoch(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	socket := cfg.ValidatorDetails.SocketPath
	// Attempt to load connection with consensus client
	connection, co := loadConsensusClient(socket)

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop code here faild to establish connection and reply
		log.Printf("Failed to establish connection using socket: %s" +
			socket)
		return
	}

	var height int64 = consensus.HeightLatest

	// Retrieve block at specific height from consensus client
	blk, err := co.GetBlock(context.Background(), height)
	if err != nil {
		log.Printf("Error while getting block info : %v", err)
		return
	}

	if &blk == nil {
		log.Printf("Got empty block res : %v", blk)
		return
	}

	bh := blk.Height

	// Return epcoh of specific height
	epoch, err := co.GetEpoch(context.Background(), bh)
	if err != nil {
		log.Printf("Failed to retrieve Epoch of Block : %v", err)
		return
	}

	err = writeToInfluxDb(c, bp, "oasis_worker_epoch_number", map[string]string{}, map[string]interface{}{"epoch_number": epoch})
	if err != nil {
		log.Printf("Error while storing worker epoch number : %v", err)
		return
	}
	log.Printf("validator worker epoch number : %v", epoch)
}
