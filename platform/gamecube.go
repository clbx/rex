package platform

import (
	"log"
	"os"

	"github.com/clbx/rex/util"
)

// Gamecube Struct
// Contains all information about a gamecube game
// Parameters
type Gamecube struct {
	Id                          string `json:"id,omitempty"`
	IsoMagicWord                string `json:"isomagicword,omitempty"`
	IsoGraphics                 string `json:"isographics,omitempty"`
	IsoShortName                string `json:"isoshortname,omitempty"`
	IsoShortDeveloper           string `json:"isoshordeveloper,omitempty"`
	IsoLongName                 string `json:"isolongname,omitempty"`
	IsoLongDeveloperDescription string `json:"isolongdevdesc,omitempty"`
	IsoDescription              string `json:"isodesc,omitempty"`
	IGDBId                      string `json:"igdbid,omitempty"`
	Path                        string `json:"path,omitempty"`
	Size                        int    `json:"size,omitempty"`
}

// Identifies gamcube game based off of ISO.
// Parameters
// path  The path to the ISO
// returns Gamecube object of the game
func IdentifyGamecube(path string) Gamecube {

	//Load ISO
	iso, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file %s", err)
	}
	defer iso.Close()

	// val := util.Read_uint32(iso, 0x001C)
	// fmt.Printf("%s", fmt.Sprintf("%X", val))

	//Gamecube disk layout: http://hitmen.c02.at/files/yagcd/yagcd/chap13.html

	//FST Offset is 4 bytes at 0x0424
	fstOffset := util.ReadUint32(iso, 0x0424)
	//FST Size is 4 bytes at 0x0428
	//fstSize := util.ReadUint32(iso, 0x0428)
	//Number of Entries
	numOfEntries := util.ReadUint32(iso, int(fstOffset)+8)

	fntOffset := fstOffset + numOfEntries*0xC

	//fmt.Printf("%d Files Found\n", numOfEntries)

	var openingOffset uint32
	//Find opening.bnr
	for i := 0; i < int(numOfEntries); i++ {
		fileOffset := fstOffset + uint32(i*0xC)
		nameOffset := util.ReadUint32(iso, int(fileOffset)) & 0x00FFFFFF
		name := util.ReadNullTerminatedString(iso, int(fntOffset+nameOffset))
		if name == "opening.bnr" {
			openingOffset = util.ReadUint32(iso, int(fileOffset)+4)
			break
		}

	}

	if openingOffset == 0 {
		log.Fatal("Unable to find opening.bnr in GCM ISO")
	}

	//fmt.Printf("opening.bnr found at %s\n", fmt.Sprintf("%X", openingOffset))

	// Gamecube Banner has information relevant to the game:
	// 0x0000 - 0x0003 Magic Word, "BNR1" (US/JP) or "BNR2" (EU)
	// 0x0004 - 0x001f nothing
	// 0x0020 - 0x181f Graphics Data.
	// 0x1820 - 0x183f Game Name
	// 0x1840 - 0x185f Developer
	// 0x1860 - 0x189f Full Game Title
	// 0x18a0 - 0x18df Company/Developer Full Name or Description
	// 0x18e0 - 0x195f Game Descriptions

	game := Gamecube{
		IsoShortName:                util.ReadNullTerminatedString(iso, int(openingOffset+0x1820)),
		IsoShortDeveloper:           util.ReadNullTerminatedString(iso, int(openingOffset+0x1840)),
		IsoLongName:                 util.ReadNullTerminatedString(iso, int(openingOffset+0x1860)),
		IsoLongDeveloperDescription: util.ReadNullTerminatedString(iso, int(openingOffset+0x18A0)),
		IsoDescription:              util.ReadNullTerminatedString(iso, int(openingOffset+0x18E0)),
		Path:                        path,
	}

	return game

}

func WrapGamecube(game Gamecube) Game {
	wrapped := Game{
		//Id:       util.GetFileMD5(game.Path),
		Name:     game.IsoLongName,
		Platform: "gcn",
		Path:     game.Path,
	}
	return wrapped
}
