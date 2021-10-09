# Golang-Sitescripts
This is a shopify module I'm writing in go below is a brief explanation of some of the code I have written.

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
