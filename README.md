# Aitunder

Beisenov Zhangeldy, Ali Samal

Group: SE-2208, SE-2219

**A brief discussion on project relevance/motivations**

The main problem is the problem of finding partners among university students. Because of this, there is a need to create a favorable environment for the development of students' social life, the creation of new ideas and projects. Therefore, we had the idea to create a website to search for like-minded people or events of interest to AITU students.

**Technical specifications of the project**

The project is a web application developed in Go (Golang) for user registration and authentication. The frontend is built using HTML with Bootstrap for styling, and the backend utilizes Go's HTTP package to handle server-side logic. The application interacts with a MongoDB database for user data storage. The code is organized into several components: HTML templates for user interfaces, Go files for server-side logic, and a MongoDB package for database operations.

Frontend:
The HTML files (login.html and registration.html) define the structure and styling of the user interfaces. Bootstrap is used for a responsive and visually appealing design. The login page provides a form for users to enter their email and password, while the registration page includes additional fields for the user's name. JavaScript is used for client-side form validation, ensuring that users provide valid information before submitting the forms.

Backend:
The Go code consists of two main parts: the server logic (main.go) and the MongoDB operations (mongodb.go). The server handles HTTP requests and routes them to the appropriate HTML templates based on the requested path ("/main" for registration and "/login" for login). It also includes API endpoints ("/api/signUp" and "/api/login") for handling user registration and authentication, respectively. The MongoDB package connects to a MongoDB database hosted on MongoDB Atlas and performs operations such as user insertion, deletion, retrieval, and authentication.


![main page](image.png)

![login](image-1.png)

**Analysis of the existing systems:**

We were inspired by the concept of Tinder, but its functionality is not suitable for creating a working atmosphere, so we found some solutions in Headhunter.

Tinder is an application for online dating and geosocial networks. On Tinder, users "swipe right" to like, or "swipe left" to dislike other users' profiles, which include their photos, a brief biography and a list of their interests. Tinder uses a "double sign-up" system where both users have to like each other before they can exchange messages.

Features:

**Swipe** is central to Tinder's design. The app's algorithm presents profile to users, who then swipe right to "like" potential matches and swipe left to continue their search.

**Messaging** is also a heavily utilized feature. Once a user matches with another user, they are able to exchange text messages on the app.

**Face to Face** is Tinder's video chat feature that allows users who have matched to see each other virtually. It was implemented in July 2020.

**Instagram** integration lets users view other users' Instagram profiles.

Headhunter is the largest online recruitment company developing business in Russia, Belarus, Kazakhstan. Headhunter's clients are over 515 thousand companies. The extensive database of applicants for HH contains more than 55 million resumes, and the average daily number of vacancies exceeds 933 thousand.


| **Advantages** | **Disadvantages** |
| --- | --- |
| Tinder |
| · Makes it possible to find new friends, partners, brings people closer;· Easy to use, with a clear interface, understandable;· There is double verification, which is a more secure applicationж | · The app has paid content;· Chances of being played or cheated, especially in relationships, because people's data is not checked for correctness;· The goal is only to find a person for a relationship; |
| HeadHunter |
| · Gives you the opportunity to quickly find a job;· Employers can select a suitable employee through the attached resume remotely;· Absolutely free application and website for people who are looking for a job; | · After registration, a lot of spam emails come to the site;· There are a lot of unnecessary ads and scams;· It is expensive for employers to place an ad;· The interface is not very clear; you need to study a little; |


**Example of the registration format:**

| Column | Datatype | Length | Attributes | Example |
| --- | --- | --- | --- | --- |
| User\_id | Int | 11 | Primary key, auto increment | 1 |
| Full Name | Varchar | 50 |
 | Ivanov Danya |
| Gender | Varchar | 6 |
 | Male |
| Login | Varchar | 20 |
 | Danya123 |
| Password | Varchar | 10 |
 | Qqwerty1! |
| Email | Varchar | 30 |
 | Danya@mail.ru |
| Phone | Int | 11 |
 | +7 777 777 88 99 |
| Group number | Varchar | 8 |
 | SE-2013 |

**Inputs and outputs of the system:**

| Inputs | Process |
| --- | --- |
| Fullname GenderLoginPasswordEmailPhoneGroup number | RegistrationAuthorizationFinding suitable project(s)Contacting with other |
| Outputs | Storage |
| Description of projects | users |

Our application also will consist of the following technical solutions:

- Swipe-able cards with swipe animations (JS)
- Matching algorithm
- Context API with a custom AuthProvider & useAuth hook
- Navigation including popup Modals

**Main functions, prototyping and prospects**

Functions:

- Registration and/or login
- Creating profile questionnaires
- Possibility of choice - refusal or approval (swipe cards)
- Matching users together who swipe on each other

Initial project drawing:

![](RackMultipart20240107-1-ktwpnh_html_5d480d72cfc98c68.jpg)

![](RackMultipart20240107-1-ktwpnh_html_484e2107e5690d30.jpg)

Sample design of the project (October 22):

![](RackMultipart20240107-1-ktwpnh_html_a38b260cb6b940fa.png)

![](RackMultipart20240107-1-ktwpnh_html_5c018b724618190b.png)

Features planned in the future:

- Google Sign In
- Integration with social networks
- News Feed (with different ads)
- Notifications from the university
- Link Microsoft Accounts

**Product Description**

