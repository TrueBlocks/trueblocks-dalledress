// Copyright 2016, 2024 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package types

// EXISTING_CODE
import (
	"encoding/json"
	"io"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// EXISTING_CODE

type MonitorEx struct {
	Address      base.Address `json:"address"`
	Deleted      bool         `json:"deleted"`
	EnsName      string       `json:"ensName"`
	FileSize     int64        `json:"fileSize"`
	Label        string       `json:"label"`
	LastScanned  uint32       `json:"lastScanned"`
	NRecords     int64        `json:"nRecords"`
	Name         string       `json:"name"`
	Stats        *Stats       `json:"stats"`
	Transactions []string     `json:"transactions"`
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s MonitorEx) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func (s *MonitorEx) Model(chain, format string, verbose bool, extraOpts map[string]any) Model {
	var model = map[string]any{}
	var order = []string{}

	// EXISTING_CODE
	// EXISTING_CODE

	return Model{
		Data:  model,
		Order: order,
	}
}

func (s *MonitorEx) MarshalCache(writer io.Writer) (err error) {
	// Address
	if err = cache.WriteValue(writer, s.Address); err != nil {
		return err
	}

	// Deleted
	if err = cache.WriteValue(writer, s.Deleted); err != nil {
		return err
	}

	// EnsName
	if err = cache.WriteValue(writer, s.EnsName); err != nil {
		return err
	}

	// FileSize
	if err = cache.WriteValue(writer, s.FileSize); err != nil {
		return err
	}

	// Label
	if err = cache.WriteValue(writer, s.Label); err != nil {
		return err
	}

	// LastScanned
	if err = cache.WriteValue(writer, s.LastScanned); err != nil {
		return err
	}

	// NRecords
	if err = cache.WriteValue(writer, s.NRecords); err != nil {
		return err
	}

	// Name
	if err = cache.WriteValue(writer, s.Name); err != nil {
		return err
	}

	// Stats
	optStats := &cache.Optional[Stats]{
		Value: s.Stats,
	}
	if err = cache.WriteValue(writer, optStats); err != nil {
		return err
	}

	// Transactions
	if err = cache.WriteValue(writer, s.Transactions); err != nil {
		return err
	}

	return nil
}

func (s *MonitorEx) UnmarshalCache(vers uint64, reader io.Reader) (err error) {
	// Check for compatibility and return cache.ErrIncompatibleVersion to invalidate this item (see #3638)
	// EXISTING_CODE
	// EXISTING_CODE

	// Address
	if err = cache.ReadValue(reader, &s.Address, vers); err != nil {
		return err
	}

	// Deleted
	if err = cache.ReadValue(reader, &s.Deleted, vers); err != nil {
		return err
	}

	// EnsName
	if err = cache.ReadValue(reader, &s.EnsName, vers); err != nil {
		return err
	}

	// FileSize
	if err = cache.ReadValue(reader, &s.FileSize, vers); err != nil {
		return err
	}

	// Label
	if err = cache.ReadValue(reader, &s.Label, vers); err != nil {
		return err
	}

	// LastScanned
	if err = cache.ReadValue(reader, &s.LastScanned, vers); err != nil {
		return err
	}

	// NRecords
	if err = cache.ReadValue(reader, &s.NRecords, vers); err != nil {
		return err
	}

	// Name
	if err = cache.ReadValue(reader, &s.Name, vers); err != nil {
		return err
	}

	// Stats
	optStats := &cache.Optional[Stats]{
		Value: s.Stats,
	}
	if err = cache.ReadValue(reader, optStats, vers); err != nil {
		return err
	}
	s.Stats = optStats.Get()

	// Transactions
	s.Transactions = make([]string, 0)
	if err = cache.ReadValue(reader, &s.Transactions, vers); err != nil {
		return err
	}

	s.FinishUnmarshal()

	return nil
}

// FinishUnmarshal is used by the cache. It may be unused depending on auto-code-gen
func (s *MonitorEx) FinishUnmarshal() {
	// EXISTING_CODE
	// EXISTING_CODE
}

// EXISTING_CODE
func NewMonitorEx(namesMap map[base.Address]NameEx, m *coreTypes.Monitor) MonitorEx {
	return MonitorEx{
		Address:     m.Address,
		Name:        namesMap[m.Address].Name.Name,
		Deleted:     m.Deleted,
		FileSize:    m.FileSize,
		LastScanned: m.LastScanned,
		NRecords:    m.NRecords,
	}
}

// EXISTING_CODE
