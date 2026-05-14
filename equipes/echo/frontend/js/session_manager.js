const SessionManager = {
    saveSession(userData) {
        const id = userData.data ? userData.data.id : (userData.user_id || userData.id);
        if (id) {
            sessionStorage.setItem('user_id', id);
            // Salva o email que vem do backend
            const email = userData.data ? userData.data.email : userData.email;
            sessionStorage.setItem('user_email', email || 'Usuário');
        }
    },
    getUserId: () => sessionStorage.getItem('user_id'),
    getUserEmail: () => sessionStorage.getItem('user_email'),
    clear() {
        sessionStorage.clear();
        window.location.replace('login.html');
    }
};