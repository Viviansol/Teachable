document.getElementById('searchForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const courseName = document.getElementById('courseName').value;
    if (!courseName) {
        alert('Please enter a course name');
        return;
    }

    fetch(`/search?course_name=${courseName}`)
        .then(response => response.json())
        .then(data => {
            if (!data.course) {
                alert('Course not found');
                return;
            }

            // Update course details table
            const courseTableBody = document.getElementById('courseTable').querySelector('tbody');
            courseTableBody.innerHTML = `
                <tr>
                    <td>${data.course.id}</td>
                    <td>${data.course.name}</td>
                    <td>${data.course.heading}</td>
                    <td>${data.course.is_published}</td>
                    <td><img src="${data.course.image_url}" alt="${data.course.name}" width="100"></td>
                </tr>
            `;

            // Update enrolled students table
            const studentsTableBody = document.getElementById('studentsTable').querySelector('tbody');
            studentsTableBody.innerHTML = '';
            data.enrollments.forEach(enrollment => {
                studentsTableBody.innerHTML += `
                    <tr>
                        <td>${enrollment.user_id}</td>
                        <td>${enrollment.enrolled_at}</td>
                        <td>${enrollment.percent_complete}%</td>
                        <td>${enrollment.user_name}</td>
                        <td>${enrollment.user_email}</td>
                    </tr>
                `;
            });
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred while fetching the data');
        });
});
