// Aguarda o carregamento do DOM para garantir que o botão exista
document.addEventListener('DOMContentLoaded', () => {
   // Seleciona o botão de logout pelo ID definido na estrutura HTML
    const btnLogout = document.getElementById('logout-btn');

    if (btnLogout) {
        btnLogout.addEventListener('click', (e) => {
            // Impede que o link recarregue a página ou siga o href="#"
            e.preventDefault();

            // Exibe uma confirmação simples
            if (confirm('Deseja realmente sair da conta?')) {
                // Chama a função de limpeza da session_manager.js
                SessionManager.clear();
            }
        });
    }
});