# Golang-Sitescripts
This is a series of golang modules I'm writing, below is a brief explanation of some of the code I have written.

<b> Shopify </b>
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

<p>
<b> Looping requests </b>
I wanted to be able to easily loop over requests when they fail while also having the ability to easily move requests around to play with website-flow order.
I decided to move away from the boolean true/false concept. Though I'm still technically setting a boolean with these new changes I will change this in future to give more control over handling specific parts of the system. I'll probably create my own type structure to represent this.
It has become clear that design between shopify sites despite using the same framework, can have largely different approaches this is especially true with addToCart.
<details>
  <summary>Start task helper</summary>
  
 ![Alt text](images/startTask.png?raw=true "example")
  
</details>
  
  
<details>
  <summary>Start task helper</summary>
  
 ![Alt text](images/startTaskImplement.png?raw=true "example")
  
</details>
<p>
<b>
Tested with the following shopify<br>
 </b>
https://shop.doverstreetmarket.com<br>
https://goodhoodstore.com<br>
https://limitededt.com<br>
