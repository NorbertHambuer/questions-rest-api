## Description

Rest API that allows you to create, update, delete and list questions.

### Questions

Questions have a simple structure. Each question has a body that defines what the candidate for a job position is supposed to answer. Then there are two or more options that the candidate can choose from. Each option has a body as well and a boolean attribute that defines whether the option is correct. At least one of the options is correct. Below is a JSON representation of a sample question.

```json
{
  "body": "Where does the sun set?",
  "options": [
    {
      "body": "East",
      "correct": false
    },
    {
      "body": "West",
      "correct": true
    }
  ]
}
```

### Endpoints

- POST /question - Creates a new question in the database and then returns it in the response
- PUT /question/{id} - Updates an existing question and returns the updated question in the response
- DELETE /question/{id} - Deletes an existing question
- GET /questions - Returns a list of all questions in the database
- GET /docs - Loads the OpenApi documentation