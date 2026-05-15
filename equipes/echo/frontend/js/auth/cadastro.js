
function checarSessao() {
    if (SessionManager.getUserId()) {
        window.location.replace('tarefas.html');
    }
}
checarSessao();
window.addEventListener('pageshow', checarSessao);

window.onload = () => {
    const form = document.getElementById('form-cadastro');

    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const email = form.email.value.trim();
            const password = form.password.value;

            try {

                await api.register(email, password);
                
                alert('Conta criada com sucesso!');
                window.location.replace('login.html');
            } catch (error) {
                alert(error.message);
            }
        });
    }
};