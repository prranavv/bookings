document.getElementById("colorButton").addEventListener("click",function(){
                    
                    let html=
                    `<form id="check-availabilty-form action="" method="post" novalidate class="needs-validation">
                          <div class="row" id="reseravation-dates-modal">
                            <div class="col">
                              <input disabled required class="form-control" type="text" id="start" name="start" placeholder="Arrival">
                            </div>
                            <div class="col">
                              <input disabled required class="form-control" type="text" id="end" name="end" placeholder="Departure">
                            </div>
                          </div>
                    </form>`
                    //notify("this is my message","warning")
                    //notifyModal("title","hello world","success","My ass")
                    attention.custom({msg:html,title:"Choose your dates"})
                  })