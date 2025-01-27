package googleplay

import (
   "fmt"
   "os"
   "strconv"
   "testing"
   "time"
)

var apps = []app_type{
   {"2022-04-28 00:00:00 +0000 UTC",2,"com.miHoYo.GenshinImpact"},
   {"2022-05-11 00:00:00 +0000 UTC",1,"com.supercell.brawlstars"},
   {"2022-05-12 00:00:00 +0000 UTC",0,"com.clearchannel.iheartradio.controller"},
   {"2022-05-23 00:00:00 +0000 UTC",0,"kr.sira.metal"},
   {"2022-05-23 00:00:00 +0000 UTC",2,"com.kakaogames.twodin"},
   {"2022-05-30 00:00:00 +0000 UTC",1,"com.madhead.tos.zh"},
   {"2022-05-31 00:00:00 +0000 UTC",1,"com.xiaomi.smarthome"},
   {"2022-06-02 00:00:00 +0000 UTC",0,"org.thoughtcrime.securesms"},
   {"2022-06-02 00:00:00 +0000 UTC",1,"com.binance.dev"},
   {"2022-06-08 00:00:00 +0000 UTC",1,"com.sygic.aura"},
   {"2022-06-12 00:00:00 +0000 UTC",0,"br.com.rodrigokolb.realdrum"},
   {"2022-06-13 00:00:00 +0000 UTC",0,"com.app.xt"},
   {"2022-06-13 00:00:00 +0000 UTC",0,"com.google.android.youtube"},
   {"2022-06-13 00:00:00 +0000 UTC",0,"com.instagram.android"},
   {"2022-06-13 00:00:00 +0000 UTC",1,"com.axis.drawingdesk.v3"},
   {"2022-06-14 00:00:00 +0000 UTC",0,"com.pinterest"},
   {"2023-02-01",1,"com.miui.weather2"},
   {"2023-02-20",0,"org.videolan.vlc"},
   {"2023-03-15",0,"com.amctve.amcfullepisodes"},
   {"2023-03-22",0,"app.source.getcontact"},
   {"2023-04-11",0,"com.google.android.apps.walletnfcrel"},
}

func Test_Details(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   var head Header
   head.Read_Auth(home + "/Documents/googleplay.txt")
   head.Auth.Exchange()
   for _, app := range apps {
      platform := Platforms[app.platform]
      head.Read_Device(home + "/Documents/" + platform + ".bin")
      d, err := head.Details(app.doc)
      if err != nil {
         t.Fatal(err)
      }
      if _, err := d.Version(); err != nil {
         t.Fatal(err)
      }
      if _, err := d.Version_Code(); err != nil {
         t.Fatal(err)
      }
      if _, err := d.Title(); err != nil {
         t.Fatal(err)
      }
      if _, err := d.Installation_Size(); err != nil {
         t.Fatal(err)
      }
      if _, err := d.Num_Downloads(); err != nil {
         t.Fatal(err)
      }
      if _, err := d.Currency_Code(); err != nil {
         t.Fatal(err)
      }
      raw_date, err := d.Upload_Date()
      if err != nil {
         t.Fatal(err)
      }
      date, err := time.Parse("Jan 2, 2006", raw_date)
      if err != nil {
         t.Fatal(err)
      }
      app.date = date.String()
      fmt.Print(app, ",\n")
      time.Sleep(99 * time.Millisecond)
   }
}

func (a app_type) String() string {
   var b []byte
   b = append(b, '{')
   b = strconv.AppendQuote(b, a.date)
   b = append(b, ',')
   b = strconv.AppendInt(b, a.platform, 10)
   b = append(b, ',')
   b = strconv.AppendQuote(b, a.doc)
   b = append(b, '}')
   return string(b)
}

type app_type struct {
   date string
   platform int64 // X-DFE-Device-ID
   doc string
}
