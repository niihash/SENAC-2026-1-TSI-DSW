const SessionManager = {
    // Salva os dados após o login bem-sucedido
    saveSession(userData) {
        if (userData && userData.id) {
            sessionStorage.setItem('user_id', userData.id);
            sessionStorage.setItem('user_email', userData.email);
        }
    },

    // Retorna o ID  usar nos fetchs de Task
    getUserId() {
        return sessionStorage.getItem('user_id');
    },

    // Limpa tudo e desloga 
    clear() {
        sessionStorage.clear();
        window.location.href = 'login.html';
    }
};