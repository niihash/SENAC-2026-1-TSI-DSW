function checarSessao() {
    const logado = SessionManager.getUserId();
    const path = window.location.pathname;

    if (logado && (path === '/' || path === '' || path.endsWith('login.html'))) {
        window.location.replace('tarefas.html');
    }
}
checarSessao();
window.addEventListener('pageshow', checarSessao);

window.onload = () => {

    const form = document.getElementById('form-login');

    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = new FormData(form);
            const email = formData.get('email');
            const password = formData.get('password');

            try {
                const userData = await api.login(email, password);
                
                if (userData && (userData.user_id || userData.ID || (userData.data && userData.data.id))) {
                    // Salva a sessão antes de pular de página
                    SessionManager.saveSession(userData);
                    window.location.replace('tarefas.html');
                } else {
                    alert("Usuário validado, mas ID não encontrado.");
                }
            } catch (error) {
                alert("Erro ao logar: " + error.message);
            }
        });
    }
};