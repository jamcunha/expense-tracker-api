# Expense Tracking App

## TODO: write the README

```
Users {
    name,
    password, // bcrypt2
    email,
}
```

```
Transactions {
    title
    transaction_type, // income or expense
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

```
Frequencies {
daily,
weekly,
monthly,
yearly,
}
```

```
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

```
Goals {
    title,
    amount_to_save,
    start_date,
    end_date,
    user, // add user as a FK
}
```

```
Savings {
    amount,
    user, // add user as a FK
}
```

#### Extra ideas

Check if there is an API to get data from bank account

#### alksdjfa

bcrypt2 for password hashing
JWT for authentication
