package googleplay

import (
   "2a.pages.dev/rosso/http"
   "2a.pages.dev/rosso/protobuf"
   "fmt"
   "io"
   "net/url"
)

func (f File) OBB(file_type uint64) string {
   var b []byte
   if file_type >= 1 {
      b = append(b, "patch"...)
   } else {
      b = append(b, "main"...)
   }
   b = append(b, '.')
   b = fmt.Append(b, f.Version_Code)
   b = append(b, '.')
   b = append(b, f.Package_Name...)
   b = append(b, ".obb"...)
   return string(b)
}

func (f File) APK(id string) string {
   var b []byte
   b = append(b, f.Package_Name...)
   b = append(b, '-')
   if id != "" {
      b = append(b, id...)
      b = append(b, '-')
   }
   b = fmt.Append(b, f.Version_Code)
   b = append(b, ".apk"...)
   return string(b)
}

func (h Header) Delivery(doc string, vc uint64) (*Delivery, error) {
   req := http.Get(&url.URL{
      Scheme: "https",
      Host: "play-fe.googleapis.com",
      Path: "/fdfe/delivery",
      RawQuery: fmt.Sprint("doc=", doc, "&vc=", vc),
   })
   h.Set_Agent(req.Header)
   h.Set_Auth(req.Header) // needed for single APK
   h.Set_Device(req.Header)
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   body, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   // ResponseWrapper
   response_wrapper, err := protobuf.Unmarshal(body)
   if err != nil {
      return nil, err
   }
   // .payload.deliveryResponse
   delivery_response := response_wrapper.Get(1).Get(21)
   // .status
   status, err := delivery_response.Get_Varint(1)
   if err != nil {
      return nil, err
   }
   switch status {
   case 2:
      return nil, fmt.Errorf("geo-blocking")
   case 3:
      return nil, fmt.Errorf("purchase required")
   case 5:
      return nil, fmt.Errorf("invalid version")
   }
   var del Delivery
   // .appDeliveryData
   del.Message = delivery_response.Get(2)
   return &del, nil
}

// .downloadUrl
func (d Delivery) Download_URL() (string, error) {
   return d.Get_String(3)
}

func (d Delivery) Split_Data() []Split_Data {
   var splits []Split_Data
   // .splitDeliveryData
   for _, split := range d.Get_Messages(15) {
      splits = append(splits, Split_Data{split})
   }
   return splits
}

// AndroidAppDeliveryData
type Delivery struct {
   protobuf.Message
}

func (d Delivery) Additional_File() []App_File_Metadata {
   var files []App_File_Metadata
   // .additionalFile
   for _, file := range d.Get_Messages(4) {
      files = append(files, App_File_Metadata{file})
   }
   return files
}

// SplitDeliveryData
type Split_Data struct {
   protobuf.Message
}

// .id
func (s Split_Data) ID() (string, error) {
   return s.Get_String(1)
}

// .downloadUrl
func (s Split_Data) Download_URL() (string, error) {
   return s.Get_String(5)
}

// AppFileMetadata
type App_File_Metadata struct {
   protobuf.Message
}

// .fileType
func (a App_File_Metadata) File_Type() (uint64, error) {
   return a.Get_Varint(1)
}

// .downloadUrl
func (a App_File_Metadata) Download_URL() (string, error) {
   return a.Get_String(4)
}

type File struct {
   Package_Name string
   Version_Code uint64
}
