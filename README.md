# Expense Tracking App

## TODO: write the README

TODO: write tests for API instead of testing manually with curl

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
    category, // add an FK or let user have it's own categories
    user, // add user as a FK
}
```

```
Categories {
    Grocery,
    Rent,
    Utilities,
    Entertainment,
    Other,
}

// explore this way
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

Check if there is an API to get data from bank account
