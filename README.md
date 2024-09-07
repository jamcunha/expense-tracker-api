# Expense Tracking App

## TODO: write the README

TODO: write tests for API instead of testing manually with curl

TODO: rewrite all deletes and updates to also use userID to make sure the user is deleting their own data

TODO: add on cascade delete to all tables

TODO: uniformize the errors log (and what to send to the client)

```
Users {
    name,
    password,
    email,
}
```

```
Expenses {
    title
    amount,
    date,
    category, // add an FK
    user, // add user as a FK
}
```

```
Categories {
    name,
    user, // add user as FK
}
```

```
Budgets {
    amount,
    start_date,
    end_date,
    user, // add user as a FK
    category, // add category as a FK
}
```

Below are ideas for future features

```
Frequencies {
    daily,
    weekly,
    monthly,
    yearly,
}
```

```
// maybe add a isPayed tracker
RecurringExpenses {
    title,
    amount,
    start_date,
    end_date,
    frequency, // add frequency as a FK
    user, // add user as a FK
    category, // add category as a FK
}
```

#### Extra ideas

Check if there is an API to get data from bank accounts
