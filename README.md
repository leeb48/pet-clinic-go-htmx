# Pet Clinic

## Purpose
- Create an applications that allows pet owners to schedule appt with veternarians

## Features
- Owner CRUD
 - Pet CRUD that belongs to an owner
- Veternarian CRUD
- Appointments can be made to a veternarian on a certain time frame
- Owner list view and search feature
- Veterinarian schedule view
 - https://fullcalendar.io/docs/initialize-globals use some kind of calendar to display schedule


## Data Models
- Owner
    - Name
    - Addr
    - City
    - State
    - phone
    - email
    - pets
- Pets
    - birthdate
    - type
    - name
    - visits
- Visits
    - veterinarian
    - pet
    - DateTime
    - Description
- Veterinarian
    - Name
    - visits