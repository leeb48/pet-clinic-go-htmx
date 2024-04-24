function renderCalendar(visits = []) {
    const events = []

    for (const visit of visits) {
        const appointment = new Date(visit.appointment)

        const startTime = appointment.toISOString()
        const endTime = new Date(appointment.setHours(appointment.getHours() + visit.duration)).toISOString()

        events.push({
            title: visit.petName,
            start: startTime,
            end: endTime,
            visitId: visit.id
        })
    }


    const calendarEl = document.getElementById('calendar');
    const calendar = new FullCalendar.Calendar(calendarEl, {
        initialView: 'timeGridWeek',
        headerToolbar: {
            left: 'prev,next',
            center: 'title',
            right: 'timeGridWeek,timeGridDay' // user can switch between the two
        },
        nowIndicator: true,
        events,
        eventClick: function (info) {
            const visitId = info.event._def.extendedProps.visitId
            window.location.href = `/visit/detail/${visitId}`
        }
    });
    calendar.render();
}
