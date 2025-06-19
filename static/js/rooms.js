// Configuration for room status updates
const STATUS_UPDATE_INTERVAL_IN_MS = 10000;
let statusUpdateTimer = null;
let roomStatusData = null;
let previousRoomData = null;

document.addEventListener("DOMContentLoaded", () => {
    focusInput();
    initRoomStatusUpdates();
    setupEventDelegation();
    
    // Add auth-user class to body if authenticated
    if (document.body.dataset.authenticated === "true") {
        document.body.classList.add("auth-user");
    }
});

// Initialize room status update system using Page Visibility API
function initRoomStatusUpdates() {
    // Initial status fetch
    fetchRoomStatus();
    
    // Set up visibility change detection
    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    // Start regular polling when visible
    if (document.visibilityState === 'visible') {
        startStatusPolling();
    }
}

function handleVisibilityChange() {
    if (document.visibilityState === 'visible') {
        // Tab became visible - fetch immediately and start polling
        fetchRoomStatus();
        startStatusPolling();
    } else {
        // Tab hidden - stop polling to save resources
        stopStatusPolling();
        window.setTimeout(() => {
            fetchRoomStatus();
        }, STATUS_UPDATE_INTERVAL_IN_MS); // One last time, to reflect having potentially joined a room
    }
}

function startStatusPolling() {
    if (!statusUpdateTimer) {
        statusUpdateTimer = setInterval(fetchRoomStatus, STATUS_UPDATE_INTERVAL_IN_MS);
    }
}

function stopStatusPolling() {
    if (statusUpdateTimer) {
        clearInterval(statusUpdateTimer);
        statusUpdateTimer = null;
    }
}

function fetchRoomStatus() {
    // Show update indicator
    const indicator = document.getElementById('update-indicator');
    indicator.classList.add('updating');
    
    fetch('/api/user/rooms')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Only update UI if data has changed
            if (!deepEqual(previousRoomData, data)) {
                updateRoomUI(data);
            }
        })
        .catch(error => {
            console.error('Error fetching room status:', error);
        })
        .finally(() => {
            // Hide indicator after a short delay
            setTimeout(() => {
                indicator.classList.remove('updating');
            }, 500);
        });
}

function updateRoomUI(data) {
    // Use a Turbolinks-style approach to update just the room list
    fetch('/')
        .then(response => response.text())
        .then(html => {
            // Create a temporary element to parse the HTML
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            // Get the room list from the fetched page
            const newRoomList = doc.querySelector('.room-list');
            
            // Preserve the current filter value
            const currentFilterValue = document.getElementById('filter-input').value;
            
            // Replace the current room list with the new one
            const currentRoomList = document.querySelector('.room-list');
            currentRoomList.innerHTML = newRoomList.innerHTML;
            
            // Reapply filtering if needed
            if (currentFilterValue) {
                filterRooms();
            }
            
            // Remove all loading states from star buttons after UI update
            document.querySelectorAll('.room-item__star.loading').forEach(button => {
                button.classList.remove('loading');
            });

            previousRoomData = JSON.parse(JSON.stringify(data)); // Deep copy
            roomStatusData = data;
        })
        .catch(error => {
            console.error('Error updating rooms UI:', error);
            
            // Remove all loading states on error
            document.querySelectorAll('.room-item__star.loading').forEach(button => {
                button.classList.remove('loading');
            });
        });
}

function setupEventDelegation() {
    // Use event delegation for copy button clicks
    document.addEventListener('click', function(event) {
        // Handle copy button clicks
        const copyButton = event.target.closest('.room-item__copy-link');
        if (copyButton) {
            const slug = copyButton.closest('.room-item').querySelector('.room-item__slug').textContent.trim();
            copyToClipboard(event, slug);
        }
        
        // Handle star button clicks
        const starButton = event.target.closest('.room-item__star');
        if (starButton) {
            event.preventDefault();
            event.stopPropagation();
            toggleStarRoom(starButton);
        }
        
        // Handle filter reset button clicks
        if (event.target.closest('#filter-reset-btn')) {
            resetFilter();
        }
    });
    
    // Handle filter input
    document.addEventListener('input', function(event) {
        if (event.target.id === 'filter-input') {
            filterRooms();
        }
    });
    
    // Handle enter key for launching rooms
    document.addEventListener('keydown', function(event) {
        if (event.target.id === 'filter-input' && event.key === 'Enter') {
            launchRoom(event);
        }
    });
}

// Function to toggle star status of a room
function toggleStarRoom(starButton) {
    const roomItem = starButton.closest('.room-item');
    const roomId = roomItem.dataset.id;
    
    if (!roomId) {
        console.error('Room ID not found');
        return;
    }
    
    // Prevent double-clicks by checking if already loading
    if (starButton.classList.contains('loading')) {
        return;
    }
    
    // Add loading state (pulsating gray filled star)
    starButton.classList.add('loading');
    
    // Save current state to revert to on error
    const wasStarred = roomItem.classList.contains('starred');
    
    // Call the API to toggle star status
    fetch(`/api/user/rooms/${roomId}/toggle-star`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        // Refresh the room list to update sorting,
        // but keep the loading animation until the refresh is complete
        fetchRoomStatus();
    })
    .catch(error => {
        console.error('Error toggling star status:', error);
        
        // Remove loading state immediately on error
        starButton.classList.remove('loading');
    });
}

function focusInput() {
    document.getElementById("filter-input").focus();
}

function copyToClipboard(event, slug) {
    event.preventDefault();
    event.stopPropagation();
    const link = `${window.location.origin}/${slug}`;
    navigator.clipboard.writeText(link);
}

function filterRooms() {
    const filter = document.getElementById("filter-input").value.toLowerCase();
    const rooms = document.querySelectorAll(".room-item");
    let visibleRoomsCount = 0;

    rooms.forEach(room => {
        const slug = room.querySelector(".room-item__slug").textContent.toLowerCase();
        if (slug.includes(filter)) {
            room.closest('li').style.display = "";
            visibleRoomsCount++;
        } else {
            room.closest('li').style.display = "none";
        }
    });

    // Check if room categories should be shown or hidden
    document.querySelectorAll('.room-list').forEach(list => {
        const hasVisibleRooms = Array.from(list.querySelectorAll('.room-item')).some(item => item.closest('li').style.display !== 'none');
        if (hasVisibleRooms) {
            list.style.display = '';
        } else {
            list.style.display = 'none';
        }
    });
    
    // Show "no rooms" message if all categories are hidden
    const allListsHidden = Array.from(document.querySelectorAll('.room-list')).every(list => list.style.display === 'none');
    const noRoomsMessage = document.querySelector('.no-rooms-message');
    if (noRoomsMessage) {
        noRoomsMessage.style.display = allListsHidden ? '' : 'none';
    }
}

function launchRoom(event) {
    if (event.key === "Enter") {
        const rooms = document.querySelectorAll(".room-item");
        const visibleRooms = Array.from(rooms).filter(room => room.closest('li').style.display !== "none");
        if (visibleRooms.length === 1) {
            const link = visibleRooms[0].querySelector(".room-item__link").href;
            window.open(link, "_blank");
        }
    }
}

function resetFilter() {
    document.getElementById("filter-input").value = "";
    filterRooms();
}

function deepEqual(obj1, obj2) {
    if (obj1 === obj2) return true;
    
    if (obj1 == null || obj2 == null) return obj1 === obj2;
    
    if (typeof obj1 !== 'object' || typeof obj2 !== 'object') return obj1 === obj2;
    
    const keys1 = Object.keys(obj1);
    const keys2 = Object.keys(obj2);
    
    if (keys1.length !== keys2.length) return false;
    
    for (let key of keys1) {
        if (!keys2.includes(key) || !deepEqual(obj1[key], obj2[key])) {
            return false;
        }
    }
    
    return true;
}