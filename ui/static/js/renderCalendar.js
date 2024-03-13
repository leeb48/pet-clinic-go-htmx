function renderCalendar(visits = []) {
    const events = []

    for (const visit of visits) {
        const startTime = new Date(visit.appointment)
        events.push({
            title: visit.petName,
            start: startTime,
            end: startTime.setHours(startTime.getHours() + visit.duration),
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
