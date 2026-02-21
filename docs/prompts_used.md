# AI Assistance Transparency Log
**Project:** Event Registration & Ticketing System REST API  
**Language/Framework:** Go (Golang) / Gin / PostgreSQL  

## Purpose of This Document
This document records the responsible use of AI-assisted development tools during the implementation of this capstone project, in accordance with guidelines permitted for academic usage. 

**AI tools were used strictly as an implementation assistant — not as an autonomous project builder.**

The vast majority of the implementation, including all core business reasoning, architectural design, database modeling, and concurrency handling strategies, was defined and manually coded by the developer. AI was used only to accelerate understanding of syntax, refine boilerplate, and unblock specific technical constraints.

## Development Philosophy
AI was utilized in the same manner as modern software engineers use tools such as:
- Official framework documentation (Gin, GORM, MDN)
- StackOverflow / GitHub Issues
- Interactive templating & syntax linting

Every snippet of code or syntax suggested by AI was:
1. Reviewed conceptually
2. Tested manually through Postman or the Web UI
3. Modified to fit the exact project requirements
4. Assessed for security and robustness

No code was included blindly without full understanding and ownership.

---

## Areas Where AI Assistance Was Used

### 1. Project Bootstrapping
**AI assistance was used to:**
- Identify the correct syntax to initialize a Gin + GORM project scaffold.
- Generate the initial boilerplate for connecting to a PostgreSQL database locally.
- Recall the exact commands to generate JWT tokens.

**Human decisions (Developer-driven):**
- Designing the Clean Architecture (Layered) folder structure (`cmd/`, `internal/models`, `/repositories`, `/services`, `/handlers`).
- Deciding which dependencies to use (e.g., opting to use `golang-jwt/jwt/v4` and native Vanilla JS over heavy external frameworks).
- Creating the relational logic connecting `Users`, `Events`, `Registrations`, and `Waitlists`.

### 2. Authentication & Data Hardening
**AI assistance was used to:**
- Pull the correct Go package (`golang.org/x/crypto/bcrypt`) for secure password hashing.
- Understand how to structure Custom Gin Middleware to extract Bearer tokens.

**Human decisions (Developer-driven):**
- Designing the Role-Based Access Control (RBAC) model allowing `AUDIENCE`, `ORGANIZER`, and `ADMIN`.
- Ensuring the hashed password is never exposed in any JSON API response.
- Designing the JWT payload to store the `user_id` mapped to UUIDs for strict schema integrity.

### 3. Concurrency Strategy (The Core Problem)
**AI assistance was used to:**
- Determine how to translate a native SQL `FOR UPDATE` lock into GORM's specific syntax (`tx.Clauses(clause.Locking{Strength: "UPDATE"})`).

**Human design decisions (Critical):**
- Recognizing the problem of race conditions in ticketing logic and discarding *Optimistic Locking* (version columns/retries) as inefficient for high contention.
- Selecting **Pessimistic Row-Level Locking** within an explicitly nested database transaction to force the database engine to queue requests serially.
- Architecting the exact logical flow inside the lock:
  1. Check `seats_remaining`.
  2. If > 0, decrement and create Registration.
  3. If = 0, shunt the user to the Waitlist table and calculate their position sequence.
  4. Ensure all logic happens before the transaction commits and the row lock releases.

### 4. Admin Simulation Tool
**AI assistance was used to:**
- Recall Go concurrency syntax, specifically exactly how `sync.WaitGroup` and `sync/atomic` counters are instantiated and safely incremented.

**Human decisions (Developer-driven):**
- Devising the idea of a live simulation endpoint to artificially trigger a DDOS-style ticket booking event to physically prove the Pessimistic Locking works.
- Managing database uniqueness constraints to ensure the dummy simulation users don't trigger SQL errors on consecutive runs (using randomized UUID generation for sample emails).

### 5. Frontend UI Development
**AI assistance was used to:**
- Generate modern CSS glassmorphism styles and color palettes.
- Look up exact Vanilla Javascript DOM manipulation syntax (e.g., `document.querySelectorAll` handling).

**Human decisions (Developer-driven):**
- Choosing to ship the entire Full Stack natively by serving static `.html` files straight out of the Go binary to bypass CORS complexity.
- Mapping out the user experience (Dashboard vs Organizer Panel vs Simulation Modal).

---

## What AI Was NOT Used For
AI was **not used** to:
- "Build the app for me."
- Choose the concurrency modeling strategy.
- Determine the REST API routes and structure.
- Write the final comprehensive `README.md` and documentation representing my understanding.
- Do any testing strategy or validation.

All system-level reasoning was strictly developer-driven.

---

## Example Harmless AI Prompts Used
To transparently demonstrate the level of assistance leveraged, here are some exact phrasing examples of prompts given to AI:

* *"What is the correct GORM syntax to apply a row-level lock (SELECT FOR UPDATE) inside a transaction block?"*
* *"Is there a quick way in Vanilla Javascript to dynamically show and hide a CSS class using a toggle?"*
* *"In Go, how do you safely increment an integer from inside multiple goroutines at the same time so the count doesn't get corrupted?"*
* *"How can I serve static HTML and CSS files directly from my Gin HTTP router so I don't need a separate frontend server?"*

---

## Statement of Responsibility
This project is a complete reflection of the developer's understanding of:
- REST API layer separation.
- Relational schema modeling and transaction integrity.
- Pessimistic vs Optimistic Database locking mechanisms.
- Go concurrency paradigms.
- JWT Security paradigms.

AI tools were leveraged successfully as a coding accelerator, purely to cut down on documentation research time. All implemented features, lines of code, and architectural design choices have been verified manually, extensively tested, and can be fully explained and defended during evaluation or viva defense.
