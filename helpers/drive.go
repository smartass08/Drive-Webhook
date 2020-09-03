package helpers
import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var wg sync.WaitGroup


type GoogleDriveClient struct{
	RootId string
	GDRIVE_DIR_MIMETYPE string
	TokenFile string
	CredentialFile string
	DriveSrv *drive.Service
	SA_ARR []string
	SA_DIR_BASE_PATH string
	SA_INDEX int
	TransferredSize int64
	DriveServices []*drive.Service
}

func(G *GoogleDriveClient) Init(rootId string){
	if rootId == ""{
		rootId = "root"
	}
	G.RootId = rootId
	G.GDRIVE_DIR_MIMETYPE = "application/vnd.google-apps.folder"
	G.TokenFile = "token.json"
	G.SA_DIR_BASE_PATH = "accounts"
	G.SA_INDEX = 0
	G.TransferredSize = 0
	files, err := ioutil.ReadDir(G.SA_DIR_BASE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range files {
		file_path := path.Join(G.SA_DIR_BASE_PATH,v.Name())
		G.SA_ARR  = append(G.SA_ARR,file_path)
	}
}

func(G *GoogleDriveClient) AuthorizeAllSa(){
	for G.SA_INDEX <= length{
		G.Authorize()
		G.SA_INDEX += 1
	}
	G.SA_INDEX = 0
}

func(G *GoogleDriveClient) SwitchServiceAccount(){
	G.SA_INDEX += 1
	G.Authorize()
}


func(G *GoogleDriveClient) Authorize(){
	b, err := ioutil.ReadFile(G.SA_ARR[G.SA_INDEX])
	config ,err := google.JWTConfigFromJSON(b,drive.DriveScope)
	if err != nil {
		log.Fatalf("failed to get JWT from JSON: %v.", err)
	}
	client := config.Client(context.Background())
	srv, err := drive.New(client)
	if err != nil {
			log.Fatal(err)
	}
	G.DriveSrv = srv
	G.DriveServices = append(G.DriveServices,srv)
}


func(G *GoogleDriveClient) CreateDir(name string, parentId string) (string, error) {
	d := &drive.File{
		Name:     name,
		MimeType: G.GDRIVE_DIR_MIMETYPE,
		Parents:  []string{parentId},
	}
	file, err := G.DriveSrv.Files.Create(d).SupportsAllDrives(true).Do()
	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		if strings.Contains(err.Error(),"User rate") {
			if G.SA_INDEX == len(G.SA_ARR)-1 {
				G.SA_INDEX = 0
			}
			G.SA_INDEX += 1
			G.DriveSrv = G.DriveServices[G.SA_INDEX]
			_, _ = G.CreateDir(name, parentId)
		}
		return "", err
	}
	fmt.Println("Created G-Drive Folder: ",file.Id)
	return file.Id, nil
}


func(G *GoogleDriveClient) GetFileMetadata(fileId string) *drive.File{
	fmt.Printf("Getting Metadata of : %s\n",fileId)
	file,err:= G.DriveSrv.Files.Get(fileId).Fields("name,mimeType,size,id,parents").SupportsAllDrives(true).Do()
	if err != nil{
		log.Fatal(err)
	}
	return file
}

func(G *GoogleDriveClient) CopyFile(fileId string, destId string) (*drive.File, error){
	f := &drive.File{
		Parents: []string{destId},
	}
	r, err := G.DriveSrv.Files.Copy(fileId, f).Do()
	if err != nil {
		log.Printf("Error while copying files %s\n", err)
		if strings.Contains(err.Error(),"User rate") {
			if G.SA_INDEX == len(G.SA_ARR)-1 {
				G.SA_INDEX = 0
			}
			G.SA_INDEX += 1
			G.DriveSrv = G.DriveServices[G.SA_INDEX]
			_, _ = G.CopyFile(fileId, destId)
		}
		return nil, err
	}
	return r,nil
}

func(G *GoogleDriveClient) CloneSa(sourceId string,parentId string){
	file := G.GetFileMetadata(sourceId)
	file,err := G.CopyFile(file.Id,parentId)
	if err != nil {
		log.Println(err)
	} else {
		wg.Wait()
		fmt.Printf("File Id: %s\n",file.Id)
	}
}

 func (G *GoogleDriveClient) DeleteFile(ID string) {
	 _ = G.DriveSrv.Files.Delete(ID).SupportsAllDrives(true).Do()
 }


func (G *GoogleDriveClient) MoveFile(ptid string, dsid string){
	file := G.GetFileMetadata(ptid)
	//Insert parent id of the same file from DB of backup drive into the variable ptid here
	///
	//After getting the the parent id
	f := &drive.File{}
	_, err := G.DriveSrv.Files.Update(ptid, f).AddParents(dsid).RemoveParents(ptid).Do()
	if err != nil{
		fmt.Printf("An error occurred while copying, switching sa: %v\n", err)
		if strings.Contains(err.Error(),"User rate") {
			if G.SA_INDEX == len(G.SA_ARR)-1 {
				G.SA_INDEX = 0
			}
			G.SA_INDEX += 1
			G.DriveSrv = G.DriveServices[G.SA_INDEX]
			G.MoveFile(file.Id, dsid)
		}
		log.Printf("Error occured while moving file %s\n", err)
		return
	}
	log.Println("Moving the file done")
}