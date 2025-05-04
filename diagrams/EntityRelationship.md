# Entity Relationships

```mermaid
erDiagram
    ORGANIZATION {
        string name
    }

    USER {
        string username
        string address
    }

    APPLICATION {
        string name
    }

    PROJECT {
        string name
    }

    PREFERENCE {
        string scope
        string key
        string value
    }

    ORGANIZATION ||--o{ APPLICATION : "creates"
    APPLICATION ||--o{ USER : "used by"
    USER ||--o{ PROJECT : "creates"

    ORGANIZATION ||--o{ PREFERENCE : "defines default"
    APPLICATION ||--o{ PREFERENCE : "customizes"
    USER ||--o{ PREFERENCE : "personalizes"
    PROJECT ||--o{ PREFERENCE : "overrides"
```
