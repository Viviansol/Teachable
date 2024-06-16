document.getElementById('searchForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const courseName = document.getElementById('courseName').value.trim();

    fetch(`/search?course_name=${courseName}`)
        .then(response => response.json())
        .then(data => {
            const resultsContainer = document.getElementById('results');
            resultsContainer.innerHTML = ''; // Clear previous results

            data.forEach(entry => {
                const course = entry.course;
                const enrollments = entry.enrollments;

                const courseTable = document.createElement('table');
                courseTable.innerHTML = `
                    <caption>${course.name}</caption>
                    <thead>
                        <tr>
                            <th>Course ID</th>
                            <th>Course Name</th>
                            <th>Description</th>
                            <th>Heading</th>
                            <th>Is Published</th>
                            <th>Image</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>${course.id}</td>
                            <td>${course.name}</td>
                            <td>${course.description}</td>
                            <td>${course.heading}</td>
                            <td>${course.is_published ? 'Yes' : 'No'}</td>
                            <td><img src="${course.image_url}" alt="${course.name}" width="100"></td>
                        </tr>
                    </tbody>
                `;

                const studentsTable = document.createElement('table');
                studentsTable.innerHTML = `
                    <caption>Enrolled Students</caption>
                    <thead>
                        <tr>
                            <th>User ID</th>
                            <th>User Name</th>
                            <th>User Email</th>
                            <th>Enrolled At</th>
                            <th>Completed At</th>
                            <th>Percent Complete</th>
                            <th>Expires At</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${enrollments.map(enrollment => `
                            <tr>
                                <td>${enrollment.user_id}</td>
                                <td>${enrollment.user_name}</td>
                                <td>${enrollment.user_email}</td>
                                <td>${enrollment.enrolled_at}</td>
                                <td>${enrollment.completed_at}</td>
                                <td>${enrollment.percent_complete}%</td>
                                <td>${enrollment.expires_at}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                `;

                resultsContainer.appendChild(courseTable);
                resultsContainer.appendChild(studentsTable);
            });
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred while fetching the data');
        });
});
