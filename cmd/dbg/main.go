package main
import ("astro/internal/efem";"fmt";"os";"strconv")
func main(){
  jd,_:=strconv.ParseFloat(os.Args[1],64)
  lat,_:=strconv.ParseFloat(os.Args[2],64)
  lon,_:=strconv.ParseFloat(os.Args[3],64)
  eps:=efem.Oblicuidad(jd)
  tsl:=efem.TiempoSidereoGreenwich(jd)+lon
  c,ok:=efem.Cuspides(tsl,lat,eps)
  fmt.Printf("%.6f %.6f %v",efem.Ascendente(tsl,lat,eps),efem.MedioCielo(tsl,eps),ok)
  for _,x:=range c{fmt.Printf(" %.6f",x)}
  fmt.Println()
}
