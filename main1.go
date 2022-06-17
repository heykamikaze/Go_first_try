package main

import ("fmt"
        "net/http"
        "html/template")//пакет с доступом к выводу на стр в вебе/терминале

type User struct {
  Name string //char?
  Age uint16 // >0
  Trust_lvl int16 //can be less then 0
  Avg_g, Hpnss float64
  Skills [] string
}

func (u User)  getAllInfo() string {
  return fmt.Sprintf("Username is %s. She is %d, her trust " +
    "level is %d", u.Name, u.Age, u.Trust_lvl)

}

func (u *User) setNewName(newName string)  {
  u.Name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
  u1 := User{"Nancee", 25, 100, 4.1, 0, []string{"C", "Go", "Oil Painting"}}
  // fmt.Fprintf(w, u1.getAllInfo())
  //fmt.Fprintf(w, "<b>Main Text</b>")//bad format
  //fmt.Fprintf(w, `<h1>Main Text</h1>
//    <b>Main Text</b>`)
    tmpl, _ := template.ParseFiles("templates/home_page.html")
    tmpl.Execute(w, u1)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Contacts page")
}

func handleRequest()  {
  http.HandleFunc("/", home_page)
  http.HandleFunc("/contacts/", contacts_page)
  http.ListenAndServe(":8080", nil)
}

func main(){
  //var name User = ...
  //u1 := User{name: "Nancee", age = 25, trust_lvl: 100, avg_g: 4.1, hpnss: 0}

  handleRequest()
}
//46 24
