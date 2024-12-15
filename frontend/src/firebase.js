// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getFirestore, doc, onSnapshot } from '@firebase/firestore';
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyBH26NjgfJVlNl6pWU2zdxmwAxtm1tR2xg",
  authDomain: "englishquizappdb.firebaseapp.com",
  projectId: "englishquizappdb",
  storageBucket: "englishquizappdb.firebasestorage.app",
  messagingSenderId: "272041351077",
  appId: "1:272041351077:web:b72e24b723452db86afa4c"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const db = getFirestore(app)

export default db