// Configuration for room status updates
const STATUS_UPDATE_INTERVAL = 15000; // 15 seconds
let statusUpdateTimer = null;
let roomStatusData = null;

document.addEventListener("DOMContentLoaded", () => {
    focusInput();
    initRoomStatusUpdates();
    setupEventDelegation();
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
    }
}

function startStatusPolling() {
    if (!statusUpdateTimer) {
        statusUpdateTimer = setInterval(fetchRoomStatus, STATUS_UPDATE_INTERVAL);
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
    
    fetch('/status')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            roomStatusData = data;
            // Always use the Turbolinks approach to update UI
            updateRoomUI(data);
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

function updateRoomUI(statusData) {
    // Use a Turbolinks-style approach to update the entire room list
    fetch('/')
        .then(response => response.text())
        .then(html => {
            // Create a temporary element to parse the HTML
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            // Get the main content area from the fetched page
            const newMainContent = doc.querySelector('main');

            // Preserve the current filter value and focus state
            const currentFilterValue = document.getElementById('filter-input').value;
            const wasFilterFocused = document.activeElement === document.getElementById('filter-input');
            
            // Replace the current main content with the new one
            const currentMain = document.querySelector('main');
            currentMain.innerHTML = newMainContent.innerHTML;
            
            // Restore the filter value and reapply filtering
            const filterInput = document.getElementById('filter-input');
            filterInput.value = currentFilterValue;
            if (currentFilterValue) {
                filterRooms();
            }
            
            // Restore focus if needed
            if (wasFilterFocused) {
                filterInput.focus();
            }
        })
        .catch(error => {
            console.error('Error updating rooms UI:', error);
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