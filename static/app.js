const API_URL = "http://localhost:8080";
let currentToken = null;
let currentUser = null; // { id, role }

// DOM Elements
const authSection = document.getElementById('auth-section');
const dashboardSection = document.getElementById('dashboard-section');
const eventsGrid = document.getElementById('events-grid');
const navActions = document.getElementById('nav-actions');
const orgPanel = document.getElementById('organizer-panel');
const adminActions = document.getElementById('admin-actions');
const simModal = document.getElementById('sim-modal');

// --- Initialization ---
document.addEventListener('DOMContentLoaded', () => {
    const savedToken = localStorage.getItem('token');
    const savedRole = localStorage.getItem('role');
    const savedId = localStorage.getItem('userId');

    if (savedToken && savedRole) {
        currentToken = savedToken;
        currentUser = { role: savedRole, id: savedId };
        loadDashboard();
    }
});

// --- UI Utilities ---
function showToast(msg, type = 'success') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = msg;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}

function switchTab(e, tab) {
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.auth-form').forEach(form => form.classList.remove('active-form'));

    e.target.classList.add('active');
    if (tab === 'login') {
        document.getElementById('login-form').classList.add('active-form');
    } else {
        document.getElementById('register-form').classList.add('active-form');
    }
}

function parseJwt(token) {
    try {
        return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
        return null;
    }
}

// --- Auth ---
async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const pwd = document.getElementById('login-pwd').value;

    try {
        const res = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password: pwd })
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Login failed');

        currentToken = data.token;
        const decoded = parseJwt(currentToken);
        currentUser = { role: decoded.role, id: decoded.user_id };

        localStorage.setItem('token', currentToken);
        localStorage.setItem('role', currentUser.role);
        localStorage.setItem('userId', currentUser.id);

        showToast('Logged in successfully!');
        loadDashboard();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const name = document.getElementById('reg-name').value;
    const email = document.getElementById('reg-email').value;
    const pwd = document.getElementById('reg-pwd').value;
    const role = document.getElementById('reg-role').value;

    try {
        const res = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password: pwd, role })
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Registration failed');

        showToast('Registration successful! Please login.');
        switchTab('login');
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function logout() {
    currentToken = null;
    currentUser = null;
    localStorage.clear();
    authSection.classList.remove('hidden');
    dashboardSection.classList.add('hidden');
    navActions.innerHTML = '';
}

// --- Dashboard ---
async function loadDashboard() {
    authSection.classList.add('hidden');
    dashboardSection.classList.remove('hidden');

    // Setup Nav
    navActions.innerHTML = `
        <span style="margin-right:1rem">Role: <span style="background:var(--primary); color:white; padding:0.25rem 0.5rem; border-radius:4px; font-size:0.8rem; font-weight:bold">${currentUser.role}</span></span>
        <button onclick="logout()" class="btn btn-outline" style="padding: 0.5rem 1rem;">Logout</button>
    `;

    // Setup Roles
    if (currentUser.role === 'ORGANIZER' || currentUser.role === 'ADMIN') {
        orgPanel.classList.remove('hidden');
        loadMyEvents();
    } else {
        orgPanel.classList.add('hidden');
    }

    loadAllEvents();
}

async function fetchWithAuth(url, options = {}) {
    return fetch(`${API_URL}${url}`, {
        ...options,
        headers: {
            ...options.headers,
            'Authorization': `Bearer ${currentToken}`
        }
    });
}

// --- Events ---
async function loadAllEvents() {
    try {
        const res = await fetchWithAuth(`/events`);
        const data = await res.json();
        if (!res.ok) throw new Error("Failed to load events");

        eventsGrid.innerHTML = '';
        if (!data.events || data.events.length === 0) {
            eventsGrid.innerHTML = '<p style="color:var(--text-muted)">No published events available right now.</p>';
            return;
        }

        data.events.forEach(ev => {
            const isFull = ev.seats_remaining === 0;
            const seatClass = isFull ? 'full' : (ev.seats_remaining < 5 ? 'low' : '');

            let btnAction = `<button onclick="bookEvent('${ev.id}')" class="btn btn-primary" style="${isFull ? 'background:var(--waitlist)' : ''}">
                ${isFull ? 'Join Waitlist' : 'Book Seat'}
            </button>`;

            // Show admin simulation button
            let simulationBtn = currentUser.role === 'ADMIN' ?
                `<button onclick="openSimModal('${ev.id}', '${ev.title.replace(/'/g, "\\'")}')" class="btn btn-danger" style="margin-top:0.5rem">⚡ Run Concurrency Sim</button>` : '';

            const dateStr = new Date(ev.event_date).toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' });

            eventsGrid.innerHTML += `
                <div class="event-card">
                    <span class="event-date-badge">${dateStr}</span>
                    <h3>${ev.title}</h3>
                    <p>${ev.description}</p>
                    <div class="event-meta">
                        <div class="meta-item">📍 ${ev.location}</div>
                        <div class="meta-item">🎟️ <span class="seats-badge ${seatClass}">${ev.seats_remaining} / ${ev.capacity} Seats</span></div>
                    </div>
                    ${btnAction}
                    ${simulationBtn}
                </div>
            `;
        });
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function bookEvent(eventId) {
    try {
        const res = await fetchWithAuth(`/events/${eventId}/register`, { method: 'POST' });
        const data = await res.json();

        if (!res.ok) throw new Error(data.error || 'Failed to book');

        showToast(data.message || 'Booked successfully!', data.waitlist ? 'waitlist' : 'success');
        loadAllEvents();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

// --- Organizer ---
async function handleCreateEvent(e) {
    e.preventDefault();
    const payload = {
        title: document.getElementById('ev-title').value,
        description: document.getElementById('ev-desc').value,
        location: document.getElementById('ev-location').value,
        event_date: new Date(document.getElementById('ev-date').value).toISOString(),
        capacity: parseInt(document.getElementById('ev-capacity').value, 10)
    };

    try {
        const res = await fetchWithAuth(`/organizer/events`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        // Immediately publish it to save clicks for the demo
        await fetchWithAuth(`/organizer/events/${data.event.id}/publish`, { method: 'POST' });

        showToast('Event Created & Published!', 'success');
        e.target.reset();
        loadAllEvents();
        loadMyEvents();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function loadMyEvents() {
    try {
        const res = await fetchWithAuth(`/organizer/events`);
        const data = await res.json();

        const myGrid = document.getElementById('my-events-grid');
        myGrid.innerHTML = '';

        if (!data.events || data.events.length === 0) {
            myGrid.innerHTML = '<p>You have not created any events.</p>';
            return;
        }

        data.events.forEach(ev => {
            myGrid.innerHTML += `
                <div class="event-card">
                    <span class="event-date-badge" style="background:#dcfce7; color:#166534;">${ev.status}</span>
                    <h3>${ev.title}</h3>
                    <div class="event-meta">
                        <div class="meta-item">Capacity: ${ev.capacity}</div>
                        <div class="meta-item">Filled: ${ev.capacity - ev.seats_remaining}</div>
                    </div>
                </div>
            `;
        });
    } catch (err) {
        console.error(err);
    }
}

// --- Concurrency Simulation ---
function openSimModal(eventId, title) {
    document.getElementById('sim-event-id').value = eventId;
    document.getElementById('sim-event-title').innerText = title;

    // Reset Display
    document.getElementById('sim-results').classList.add('hidden');
    document.getElementById('res-total').innerText = '0';
    document.getElementById('res-success').innerText = '0';
    document.getElementById('res-waitlist').innerText = '0';
    document.getElementById('res-fail').innerText = '0';
    document.getElementById('res-seats').innerText = '0';

    simModal.classList.remove('hidden');
}

function closeSimModal() {
    simModal.classList.add('hidden');
}

async function handleSimulation(e) {
    e.preventDefault();
    const eventId = document.getElementById('sim-event-id').value;
    const users = document.getElementById('sim-users-count').value;

    // Show loading state
    const btn = e.target.querySelector('button');
    const originalText = btn.innerText;
    btn.innerText = 'Running Attack... ⏳';
    btn.disabled = true;

    document.getElementById('sim-results').classList.add('hidden');

    try {
        const res = await fetchWithAuth(`/admin/events/${eventId}/simulate?users=${users}`, { method: 'POST' });
        const data = await res.json();

        if (!res.ok) throw new Error(data.error || 'Simulation Failed');

        const r = data.simulation_results;

        // Display
        document.getElementById('sim-results').classList.remove('hidden');
        document.getElementById('res-total').innerText = r.total_attempted;
        document.getElementById('res-success').innerText = r.success_count;
        document.getElementById('res-waitlist').innerText = r.waitlisted_count;
        document.getElementById('res-fail').innerText = r.failed_count;
        document.getElementById('res-seats').innerText = r.final_seats_remaining;

        showToast('Simulation Complete! Check results.', 'success');
        loadAllEvents(); // Refresh seats behind
    } catch (err) {
        showToast(err.message, 'error');
    } finally {
        btn.innerText = originalText;
        btn.disabled = false;
    }
}
