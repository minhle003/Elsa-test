import React from 'react';
import CreateQuiz from '../components/CreateQuiz';

const CreateQuizPage = () => {
    const userId = localStorage.getItem('User-ID'); // Retrieve userId from localStorage

    if (!userId) {
        return <div>Please log in to create a quiz.</div>;
    }
    return (
        <div>
            <CreateQuiz userId={userId} />
        </div>
    );
};

export default CreateQuizPage;