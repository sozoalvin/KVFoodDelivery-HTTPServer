# KVFoodDelivery
Backend Prototyping System for Food Delivery Applications. Written in go/ golang with search, auto complete, suggestive search as well postal code lookups

<h2>Introduction</h2>

<p>The KV app was created in an attempt to solve teething customer service issues that have been affecting customer’s average order values, customer satisfaction and ultimate customer loyalty.
While there were alternatives that were created before; there is no single application that can integrate the features that are now available.</p>

<img src = "https://i.ibb.co/P4n3bSz/1.png">

<h2>1. Design Considerations:</h2>

<p>The CLI application should be lightweight, having the options for offline and/or online data importation. UX of the CLI application must also be user friendly to encourage early adopters to start feedback / app revisions. The main functionality of the app must also allow fast searching of certain data especially merchant information, transaction IDs as well as system queue numbers.</p>

<h2>2. Performance Considerations:</h2>

<p>Without querying online databases, databases have to be imported locally, created and stored in memory. Utilizing GO routines are mandatory as part of performance control and flow. The CLI application should also ensure that the data can be populated accurately with the help of channels confirming routine completions before the system can be utilized.
As a query app loaded with features, the features themselves must also be able to present data fast and correctly.</p>

<h2>3. Error and Panic Considerations:</h2>

<p>Error handling and panic cases are to be meticulously tested and addressed for.
Handling data, writing and reading can cause panic issues especially with databases that are heavy on slices. The most common panic that can happen in the application is due to array or slice length mismatches that occur. The implementation of a prefix trie (discussed more in the search and auto complete section) can only contain certain characters.
Characters that are recognized and stored in the arrays will cause panics. In light of that; deferred recovery functions have to be utilised in the code to mitigate such issues and recover from potential panics.</p>

<h2>Other Resources:</h2>
<p>Other resources regarding the implementaion of data structure / time complexity anaylysis can be found in the PDF and the Word Folders</p>

<p>Please feel free to write to me at sozoalvin@gmail.com if you have any questions/ feedback or if you wish to collaborate working on the application. You can also use the issues tab.</p>
