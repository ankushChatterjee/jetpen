# Jetpen
### The world's best junk email generator

This is a simple fullstack newsletter management application to learn and teach fullstack web dev.

**Features :**

1. User regsitration (With verification mail)
2. Create newsletters
3. Edit / Delete Newsletters
4. Let people subscribe/unsubscribe to the newsletter
5. Send the subscribed people email newsletters. (Each sent email from a newsletter is called a "letter")
6. Ability to save the letter as a draft.
7. Send the saved draft letter.
8. See all your past sent letters.

## Backend Architecture

In this section, we see the backend of the application. 

⚠ Note that this design is no way meant for scale. This is just to learn few aspects of a backend system not all of them. DONT PUT IT INTO PRODUCTION!

### Components of the Architecture

1. Database - PostgreSQL 🐘
2. Queueing Service - RabbitMQ 🐰
3. Microservices ⚙
    1. Email Service - Sends emails to user, reading from a queue.
    2. Newsletter Service - Handles saving, editing, deleting newsletters and letters.
    3. Sub service - Handles subscriptions and unsubscribe
    4. User service - User authentication and user profile management
4. nginx as a API Gateway
5. Backend Language - Go 💙
6. Library - gofiber ⚡

### Database DDL
Used by : `User Service`
```
create table jetpen.users(username varchar(50) primary key, email varchar(50) unique not null, name varchar(50) not null, password varchar(1024) not null, "CreatedAt" timestamp default current_timestamp);

create table jetpen.temp_users(username varchar(50) primary key, email varchar(50) unique not null, name varchar(50), password varchar(1024), "CreatedAt" timestamp default current_timestamp, token varchar(1024));
```

Used by : `Newsletter Service`
```
create table jetpen.newsletter (id uuid DEFAULT uuid_generate_v4(), name varchar(50)not null, 
				description varchar(500), 
				owner varchar(50) not null,
				"CreatedAt" timestamp default current_timestamp,
				primary KEY(id), FOREIGN key(owner) REFERENCES jetpen.users(username));
			
			
CREATE INDEX idx_jetpen_newsletter ON jetpen.newsletter ("CreatedAt", id, owner);
create table jetpen.letter(id uuid, subject varchar(200) not null, owner varchar(50) not null, content text, nid uuid not null, "CreatedAt" timestamp default current_timestamp, "isPublished" boolean, "PublishedAt" timestamp,primary key(id), FOREIGN key(owner) REFERENCES jetpen.users(username) ON DELETE CASCADE, FOREIGN key(nid) REFERENCES jetpen.newsletter(id)); 					
CREATE INDEX idx_jetpen_letter ON jetpen.letter ("CreatedAt", id);
ADD CONSTRAINT letter_nid_fkey
    FOREIGN KEY (nid)
    REFERENCES jetpen.newsletter(id)
    ON DELETE CASCADE ON UPDATE NO ACTION;

```

Used by : `Sub service`
```
create table jetpen.subscription(id uuid, email varchar(50) not null, nid uuid not null, "CreatedAt" timestamp default current_timestamp, subToken varchar(512),primary key(id), FOREIGN key(nid) REFERENCES jetpen.newsletter(id));
create index idx_jetpen_subscription_nid_email on jetpen.subscription(nid, email);
```

## Frontend
For frontend we use : 
1. React
2. React Router
3. Ant Design
4. create react app for tooling and starting template.
