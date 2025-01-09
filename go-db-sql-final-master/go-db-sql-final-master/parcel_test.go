package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const (
	ParcelStatusReg     = "registered"
	ParcelStatusShipped = "shipped"
)

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	var parsedTime time.Time
	parsedTime, err = time.Parse(time.RFC3339, "2024-12-31T23:59:59Z")
	require.NoError(t, err)

	parcel := Parcel{
		Client:    1000,
		Status:    ParcelStatusReg,
		Address:   "test",
		CreatedAt: parsedTime.Format(time.RFC3339),
	}

	id, err := store.Add(parcel)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, storedParcel)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	var parsedTime time.Time
	parsedTime, err = time.Parse(time.RFC3339, "2024-12-31T23:59:59Z")
	require.NoError(t, err)

	parcel := Parcel{
		Client:    1000,
		Status:    ParcelStatusReg,
		Address:   "test",
		CreatedAt: parsedTime.Format(time.RFC3339),
	}

	id, err := store.Add(parcel)
	require.NoError(t, err)

	err = store.SetStatus(id, ParcelStatusShipped)
	require.NoError(t, err)

	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusShipped, updatedParcel.Status)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	var parsedTime time.Time
	parsedTime, err = time.Parse(time.RFC3339, "2024-12-18T12:00:00Z")
	require.NoError(t, err)

	parcel := Parcel{
		Client:    1000,
		Status:    ParcelStatusReg,
		Address:   "test",
		CreatedAt: parsedTime.Format(time.RFC3339),
	}

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "new address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, updatedParcel.Address)
}
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	var parsedTime time.Time
	parsedTime, err = time.Parse(time.RFC3339, "2024-12-18T12:00:00Z")
	require.NoError(t, err)

	parcel1 := Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "address 1",
		CreatedAt: parsedTime.Format(time.RFC3339),
	}

	parcel2 := Parcel{
		Client:    1000,
		Status:    ParcelStatusShipped,
		Address:   "address 2",
		CreatedAt: parsedTime.Format(time.RFC3339),
	}

	_, err = store.Add(parcel1)
	require.NoError(t, err)
	_, err = store.Add(parcel2)
	require.NoError(t, err)

	parcels, err := store.GetByClient(1000)
	require.NoError(t, err)

	require.Len(t, parcels, 2)

	require.Equal(t, parcel1.Client, parcels[0].Client)
	require.Equal(t, parcel2.Client, parcels[1].Client)
	require.Equal(t, parcel1.Address, parcels[0].Address)
	require.Equal(t, parcel2.Address, parcels[1].Address)
	require.Equal(t, parcel1.Status, parcels[0].Status)
	require.Equal(t, parcel2.Status, parcels[1].Status)
}
