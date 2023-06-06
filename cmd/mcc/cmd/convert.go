package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts existing mc project to work with mcc.",
	Long:  `Converts existing mc project to work with mcc.`,
	Run:   runConvertCmd,
}

func runConvertCmd(_ *cobra.Command, args []string) {
	db := mcdb.MustConnectToDB()
	localFileStor := stor.NewGormLocalFileStor(db)
	remoteFileStor := stor.NewGormRemoteFileStor(db)
	fileStor := stor.NewGormFileStor(db)

	fmt.Println("Starting conversion...")
	err := remoteFileStor.ListPaged(func(rf *model.RemoteFile) error {
		lf, err := localFileStor.GetRemoteFileByPath(rf.Path)
		if err != nil {
			// ignore errors
			return nil
		}

		sec, dec := math.Modf(rf.MTime)
		rmTime := time.Unix(int64(sec), int64(dec*(1e9)))

		sec, dec = math.Modf(lf.MTime)
		lmTime := time.Unix(int64(sec), int64(dec*(1e9)))

		ftype := mcc.FTypeFile

		if rf.OType == "directory" {
			ftype = mcc.FTypeDirectory
		}

		f := model.File{
			RemoteID:  uint(rf.ID),
			Path:      rf.Path,
			LMTime:    lmTime,
			RMTime:    rmTime,
			LChecksum: lf.Checksum,
			RChecksum: rf.Checksum,
			FType:     ftype,
		}

		_, err = fileStor.GetFileByPath(f.Path)
		if err == nil {
			fmt.Printf("  File already converted %s, skipping...\n", f.Path)
			return nil
		}

		addedFile, err := fileStor.AddFile(f)
		if err != nil {
			fmt.Printf("  Failed to add file: %s\n", err)
			return nil
		}

		fmt.Printf("  Converted file %s...\n", addedFile.Path)
		return nil
	})

	if err != nil {
		fmt.Printf("Error converting database: %s\n", err)
	} else {
		fmt.Println("Done.")
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
