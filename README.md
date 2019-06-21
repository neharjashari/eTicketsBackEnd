# eTicketsBackEnd

Back End service for the Android application eTickets. This project was worked for the Mobile Programming course.

The app is not deployed anywhere it works only on "localhost:8000" right now.

The token provided needs to be the same format as a UUID v4.

## Instructions



### GET     /info

For this get request you get meta information about the API

    {
        "info": [
            "E-TICKETS Application",
            "Back End Service",
            "Service for Creating Events"
        ],
        "version": "v2.2.0"
    }




### POST    /event/\<token\>

You can do create events with this request by putting in the request body the following data:

    {
      "id": "",
      "title": "",
      "author": "",
      "date_created": "",
      "content": "",
      "photo": ""
    }

As a result you get the ID of the created event.



### GET     /event/\<token\>

Return the events (as an array) created by the user with the specified token in the path



### GET     /event/\<token\>/\<id\>
  
Returns the data for the event with the ID specified in the path



### GET     /event/\<token\>/tickets
  
Returns the tickets purchased from user with the specific token




There are some other API calls that can only be made by the admin and are not going to be shown here.



## Built With

Go Language



## Database

MongoDB
