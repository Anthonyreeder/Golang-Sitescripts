# Golang-Sitescripts
This is a series of golang modules I'm writing, below is a brief explanation of some of the code I have written.

<b>Shopify </b> <br>
<p>
Shopify.go is the entry point to this module. Currently set to run the front-end shopify methods.
In future i'll want to run the modules in a goroutine so i can run modules concurrently. I'll build out a task structure which loads all the required outside (UI) data which is passed into a goroutine for execution.

<details>
  <summary>Shopify package</summary>
  
 ![Alt text](images/ShopifyPackage.png?raw=true "ShopifyPackage")
  
</details>
<details>
  <summary>Shopify demo</summary>
  
 ![Alt text](images/ShopifyDemo.png?raw=true "ShopifyDemo")
  
</details>
</p>
-------------------------------------------------------------------------------------------------------------------------------------------------------
<b>Helper functions </b> <br>
<p>
These are my helper functions to keep task logic small and compact. Each module has their own helpers.go. I tried to globalise this but having that kind of relationship across packages just isn't worth the trouble. The AddHeaders method builds up a Header obj giving nicer readability without sacrificing any customisation.
The ExtractValue function will currently return a panic due to an error check on the client (in client.go) but I will change this in future when I create the GoRoutine tasks.
<details>
  <summary>Shopify helpers</summary>
  
 ![Alt text](images/ShopifyHelpers.png?raw=true "ShopifyHelpers")
  
</details>
<details>
  <summary>Shopify requests</summary>
  
 ![Alt text](images/requests.png?raw=true "requests")
  
</details>

and here is an example of a GET and POST request using the current setup.
<details>
  <summary>Shopify example</summary>
  
 ![Alt text](images/example.png?raw=true "example")
  
</details>
</p>
-------------------------------------------------------------------------------------------------------------------------------------------------------
<b> Looping requests </b><br>
<p>
I wanted to be able to easily loop over requests when they fail while also having the ability to easily move requests around to play with website-flow order.
I decided to move away from the boolean true/false concept. Though I'm still technically setting a boolean with these new changes I will change this in future to give more control over handling specific parts of the system. I'll probably create my own type structure to represent this.
It has become clear that design between shopify sites despite using the same framework, can have largely different approaches this is especially true with addToCart.
<Br>
-------------------------------------------------------------------------------------------------------------------------------------------------------
<br><b> <u>Go Routines</u> </b> <br>
I added Goroutine capability to the startTask helper function which I intend to enhance more later.
The reason I added this, certain requests we don't really care about waiting for the response. We can just submit these and then go onto to subsequent requests to save time. If the previous requests fail they will be caught later and raise a panic. In future this will force a complete task restart.
<details>
  <summary>Start task helper</summary>
  
 ![Alt text](images/startTask.png?raw=true "example")
  
</details>
  
  
<details>
  <summary>Start task implementation</summary>
  
 ![Alt text](images/startTaskImplement.png?raw=true "example")
  
</details>
<p>
<b>
Tested with the following sites<br>
 </b>

  <b>Success</b>
```  
1. atmos - https://www.atmosusa.com/ ->               100% Success
2. bdgastore - https://bdgastore.com/ ->              100% Success
3. bb branded - https://www.bbbranded.com/ ->         100% Success
4. blends US - https://www.blendsus.com/ ->           100% Success
5. coporategotem -> https://corporategotem.com ->     100% Success
6. dopefactory - https://www.dope-factory.com/ ->     100% Success
7. deadstock - https://www.deadstockofficial.com/ ->  100% Success
8. Limited Edition - https://limitededt.com/ ->       100% Success
9. goodhood - https://goodhoodstore.com/ ->           100% Success
10. sneakerboxshop.ca - https://sneakerboxshop.ca/ -> 100% success
11. Xhibition - https://www.xhibition.co/ ->          100% success
12. juicestore - https://juicestore.com/ ->           100% Success
13. Public school NY - https://www.publicschoolnyc.com/ -> 100% Success
14. saintalfred - https://www.saintalfred.com/ ->          100% Success
15. Sneaker Politics - https://sneakerpolitics.com/ ->     100% Success
16. apbstore - https://www.apbstore.com/ ->                100% Success
17. social status - https://www.socialstatuspgh.com/ ->    100% Success
18. culturekings -> https://www.culturekings.com ->        100% Success
19. feature - https://feature.com/ ->                      100% Success but could use optimising
20. Closet Inc - https://www.theclosetinc.com/ ->          100% Success but could use optimising
21. shoepalace - https://www.shoepalace.com/ ->            100% success
22. dtlr - https://www.dtlr.com/ ->                        100% success
23. ficegallery - https://www.ficegallery.com/ ->          100% success
24. just don - https://justdon.com/ ->                     100% success
25. solefiness - https://www.solefiness.com/ ->            100% success
```
  <b> Failure/In works </b>
```
smets.lu - https://smets.lu/ ->                            Paypal only
burn rubber - https://burnrubbersneakers.com/ ->           Paypal only
sole steal - https://www.solesteals.com/ ->                Paypal only
lustmexico - https://www.lustmexico.com/ ->                Paypal only
Noirfonce EU - https://www.noirfonce.eu/ ->                3DS (Confirm payment in revolut)
hanon shop - https://www.hanon-shop.com/ ->                3DS (Confirm payment in revolut)
undefeated - https://undefeated.com/ ->                    Requires login
A ma maniere - https://www.a-ma-maniere.com/ ->            Requires login
bouncewear - https://bouncewear.com/ ->                    Requires login
exclusity life - https://shop.exclucitylife.com/ ->        Requires login
jimmyjazz - https://www.jimmyjazz.com ->                   Requires login
kith - https://kith.com/ ->                                Requires login
concepts - https://cncpts.com/ ->                          Antibot/Password or offline
Funko - https://www.funko.com/ ->                          Antibot/Password or offline
haven - https://havenshop.com/ ->                          Antibot/Password or offline
Palace SB - https://www.palaceskateboards.com ->           Antibot/Password or offline
travis scott - https://www.t qravisscott.com/ ->           Antibot/Password or offline
cpfm - https://cactusplantfleamarket.com/ ->               Antibot/Password or offline
DDT - https://ddtstore.com/password ->                     Antibot/Password or offline
kawsone - https://kawsone.com/ ->                          Antibot/Password or offline
NIce Kick - https://www.nicekicks.com/ ->                  Antibot/Password or offline
ronniefieg - https://shop.ronniefieg.com/ ->               Antibot/Password or offline
Dover street market - https://london.doverstreetmarket.com/ -> Antibot/Password or offline
ovo - https://uk.octobersveryown.com/ ->                       Fails on get shipping ID for USA canada and singapour
oneblock down - https://eu.oneblockdown.it/ ->                 Fails on get shipping ID for USA canada and singapour
Packershoes - https://packershoes.com/ ->                      Fails on get shipping ID for USA canada and singapour
Suede Store - https://suede-store.com/ ->                      Fails on get shipping ID for USA canada and singapour
Trophy Room - https://www.trophyroomstore.com/ ->              Fails submitting customer information ```
