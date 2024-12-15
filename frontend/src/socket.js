import io from 'socket.io-client';

const socket = io('http://localhost:8080', {
    transports: ['websocket']
});

socket.on('connect', () => {
    console.log('Connected to server');
});

// Listen for the 'error' event
socket.on('error', (error) => {
    console.error('Error:', error);
});

socket.on('disconnect', (reason) => {
    console.log('Disconnected from server:', reason);
});


export default socket;