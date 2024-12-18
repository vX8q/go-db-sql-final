package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
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
	parcel := Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: "2024-12-18T12:00:00Z",
	}

	id, err := store.Add(parcel)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, storedParcel.Client)
	require.Equal(t, parcel.Status, storedParcel.Status)
	require.Equal(t, parcel.Address, storedParcel.Address)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}
