const API_URL = 'http://localhost:8080/api';

console.log('Auth system loading... API_URL:', API_URL);

// Получаем элементы модальных окон
const modal = document.getElementById('registerModal');
const loginModal = document.getElementById('loginModal');
const registerLink = document.getElementById('registerLink');
const closeBtns = document.querySelectorAll('.close');

// Элементы формы регистрации
const regUsername = document.getElementById('regUsername');
const regEmail = document.getElementById('regEmail');
const regPassword = document.getElementById('regPassword');
const regConfirmPassword = document.getElementById('regConfirmPassword');
const registerForm = document.getElementById('registerForm');
const passwordMatch = document.getElementById('passwordMatch');

// Элементы формы входа
const loginForm = document.querySelector('.login-form'); // старая форма в хедере
const loginModalForm = document.getElementById('loginModalForm');
const loginUsername = document.getElementById('loginUsername');
const loginPassword = document.getElementById('loginPassword');

// Кнопки переключения между окнами
const switchToLogin = document.getElementById('switchToLogin');
const switchToRegister = document.getElementById('switchToRegister');

// Элементы для валидации пароля
let passwordStrengthElement = null;
let passwordErrorsElement = null;

// Храним правила паролей
let passwordRules = null;

// Загружаем правила паролей при загрузке
async function loadPasswordRules() {
    try {
        const response = await fetch(`${API_URL}/password-rules`);
        if (response.ok) {
            passwordRules = await response.json();
            console.log('Password rules loaded:', passwordRules);
        }
    } catch (error) {
        console.warn('Could not load password rules:', error);
    }
}

// Открываем модальное окно регистрации
if (registerLink) {
    registerLink.addEventListener('click', function(event) {
        event.preventDefault();
        console.log('Opening registration modal');
        modal.style.display = 'block';
        loadPasswordRules();
        initPasswordValidation();
    });
}

// Закрываем все модальные окна
closeBtns.forEach(btn => {
    btn.addEventListener('click', function() {
        modal.style.display = 'none';
        if (loginModal) loginModal.style.display = 'none';
        resetPasswordValidation();
    });
});

// Закрываем модальные окна при клике вне их
window.addEventListener('click', function(event) {
    if (event.target === modal) {
        modal.style.display = 'none';
        resetPasswordValidation();
    }
    if (loginModal && event.target === loginModal) {
        loginModal.style.display = 'none';
    }
});

// Закрываем модальные окна при нажатии Escape
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        modal.style.display = 'none';
        if (loginModal) loginModal.style.display = 'none';
        resetPasswordValidation();
    }
});

// Переключение между окнами регистрации и входа
if (switchToLogin) {
    switchToLogin.addEventListener('click', function(event) {
        event.preventDefault();
        modal.style.display = 'none';
        loginModal.style.display = 'block';
        resetPasswordValidation();
    });
}

if (switchToRegister) {
    switchToRegister.addEventListener('click', function(event) {
        event.preventDefault();
        loginModal.style.display = 'none';
        modal.style.display = 'block';
        loadPasswordRules();
        initPasswordValidation();
    });
}

// Инициализация валидации пароля
function initPasswordValidation() {
    if (!passwordStrengthElement) {
        passwordStrengthElement = document.getElementById('passwordStrength');
    }
    if (!passwordErrorsElement) {
        passwordErrorsElement = document.getElementById('passwordErrors');
    }
    
    // Сбрасываем стили полей
    resetFieldStyles();
}

// Сбрасываем стили полей
function resetFieldStyles() {
    const inputs = [regUsername, regEmail, regPassword, regConfirmPassword];
    inputs.forEach(input => {
        if (input) {
            input.classList.remove('error-field', 'success-field');
        }
    });
}

// Сброс валидации пароля
function resetPasswordValidation() {
    if (passwordStrengthElement) {
        passwordStrengthElement.style.display = 'none';
        passwordStrengthElement.innerHTML = '';
    }
    if (passwordErrorsElement) {
        passwordErrorsElement.style.display = 'none';
        passwordErrorsElement.innerHTML = '';
        passwordErrorsElement.style.backgroundColor = '';
        passwordErrorsElement.style.color = '';
    }
    resetFieldStyles();
}

// Показываем подсказки при фокусе
if (regUsername) {
    regUsername.addEventListener('focus', function() {
        const hint = this.parentNode.querySelector('.form-hint');
        if (hint) hint.style.display = 'block';
    });
    
    regUsername.addEventListener('blur', function() {
        const hint = this.parentNode.querySelector('.form-hint');
        if (hint) hint.style.display = 'none';
        validateUsernameField(this);
    });
}

if (regEmail) {
    regEmail.addEventListener('focus', function() {
        const hint = this.parentNode.querySelector('.form-hint');
        if (hint) hint.style.display = 'block';
    });
    
    regEmail.addEventListener('blur', function() {
        const hint = this.parentNode.querySelector('.form-hint');
        if (hint) hint.style.display = 'none';
        validateEmailField(this);
    });
}

// Валидация имени пользователя
function validateUsernameField(field) {
    const value = field.value.trim();
    if (!value) return false;
    
    // Проверка длины
    if (value.length < 3) {
        field.classList.add('error-field');
        field.classList.remove('success-field');
        return false;
    }
    
    // Проверка символов
    const usernameRegex = /^[a-zA-Z0-9._-]+$/;
    if (!usernameRegex.test(value)) {
        field.classList.add('error-field');
        field.classList.remove('success-field');
        return false;
    }
    
    field.classList.remove('error-field');
    field.classList.add('success-field');
    return true;
}

// Валидация email
function validateEmailField(field) {
    const value = field.value.trim();
    if (!value) return false;
    
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) {
        field.classList.add('error-field');
        field.classList.remove('success-field');
        return false;
    }
    
    field.classList.remove('error-field');
    field.classList.add('success-field');
    return true;
}

// Показывает силу пароля
function showPasswordStrength(password) {
    if (!passwordStrengthElement) {
        initPasswordValidation();
    }
    
    if (password.length === 0) {
        passwordStrengthElement.style.display = 'none';
        return;
    }
    
    // Проверка пароля на стороне клиента
    const validation = validatePassword(password);
    
    // Определяем уровень сложности
    let strengthText = '';
    let strengthClass = '';
    let strengthPercent = 0;
    
    if (validation.score >= 80) {
        strengthText = 'Очень сильный пароль';
        strengthClass = 'strength-very-strong';
        strengthPercent = 100;
    } else if (validation.score >= 60) {
        strengthText = 'Сильный пароль';
        strengthClass = 'strength-strong';
        strengthPercent = 75;
    } else if (validation.score >= 40) {
        strengthText = 'Средний пароль';
        strengthClass = 'strength-good';
        strengthPercent = 50;
    } else if (validation.score >= 20) {
        strengthText = 'Слабый пароль';
        strengthClass = 'strength-fair';
        strengthPercent = 25;
    } else {
        strengthText = 'Очень слабый пароль';
        strengthClass = 'strength-weak';
        strengthPercent = 10;
    }
    
    passwordStrengthElement.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
            <strong style="font-size: 14px;">Сила пароля:</strong>
            <span style="font-size: 13px; color: #666;">${validation.score}/100</span>
        </div>
        <div class="password-strength-meter">
            <div class="password-strength-meter-fill ${strengthClass}" 
                 style="width: ${strengthPercent}%"></div>
        </div>
        <div style="margin-top: 8px; font-size: 13px; color: #666;">
            ${strengthText}
        </div>
    `;
    passwordStrengthElement.style.display = 'block';
    
    // Показываем ошибки, если есть
    if (validation.errors.length > 0 && passwordErrorsElement) {
        passwordErrorsElement.innerHTML = `
            <div style="color: #e74c3c; font-weight: bold; margin-bottom: 8px;">
                ⚠️ Необходимо исправить:
            </div>
            <ul style="margin: 0; padding-left: 20px; color: #e74c3c; font-size: 13px;">
                ${validation.errors.map(error => `<li>${error}</li>`).join('')}
            </ul>
        `;
        passwordErrorsElement.style.backgroundColor = '#fff5f5';
        passwordErrorsElement.style.border = '1px solid #ffcccc';
        passwordErrorsElement.style.display = 'block';
    } else if (passwordErrorsElement) {
        passwordErrorsElement.style.display = 'none';
    }
}

// Валидация пароля на стороне клиента
function validatePassword(password) {
    const errors = [];
    let score = 0;
    
    // Длина
    if (password.length < 8) {
        errors.push('Минимум 8 символов');
    } else {
        score += 20;
    }
    
    // Сложность
    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password);
    
    if (!hasUpper) errors.push('Хотя бы одна заглавная буква');
    else score += 20;
    
    if (!hasLower) errors.push('Хотя бы одна строчная буква');
    else score += 20;
    
    if (!hasNumber) errors.push('Хотя бы одна цифра');
    else score += 20;
    
    if (!hasSpecial) errors.push('Хотя бы один специальный символ');
    else score += 20;
    
    // Простые пароли
    const commonPasswords = [
        'password', '123456', '12345678', '123456789', '1234567890',
        'qwerty', 'abc123', 'password1', 'admin', 'letmein'
    ];
    
    const lowerPassword = password.toLowerCase();
    for (const common of commonPasswords) {
        if (lowerPassword.includes(common)) {
            errors.push('Слишком простой и распространенный пароль');
            score -= 30;
            break;
        }
    }
    
    // Последовательности
    if (/(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)/i.test(password) ||
        /(012|123|234|345|456|567|678|789|890)/.test(password)) {
        errors.push('Слишком простые последовательности символов');
        score -= 20;
    }
    
    // Повторяющиеся символы
    if (/(.)\1\1/.test(password)) {
        errors.push('Слишком много повторяющихся символов');
        score -= 15;
    }
    
    // Ограничение оценки
    if (score < 0) score = 0;
    if (score > 100) score = 100;
    
    return { score, errors };
}

// Проверка пароля на сервере (более точная)
async function validatePasswordOnServer(password) {
    try {
        const response = await fetch(`${API_URL}/validate-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password })
        });
        
        if (response.ok) {
            return await response.json();
        }
        return { valid: false, errors: ['Ошибка проверки пароля'] };
    } catch (error) {
        console.error('Password validation error:', error);
        return { valid: false, errors: ['Не удалось проверить пароль'] };
    }
}

// Простая проверка соединения с сервером
async function checkServerConnection() {
    try {
        console.log('Checking server connection to:', API_URL + '/health');
        const response = await fetch(API_URL + '/health');
        console.log('Server response status:', response.status);
        
        if (response.ok) {
            const data = await response.json();
            console.log('Server is healthy:', data);
            return true;
        } else {
            console.error('Server returned error:', response.status);
            return false;
        }
    } catch (error) {
        console.error('Cannot connect to server:', error);
        showNotification('❌ Не могу подключиться к серверу. Запустите бэкенд на localhost:8080', 'error');
        return false;
    }
}

// Валидация совпадения паролей в реальном времени
if (regPassword && regConfirmPassword) {
    regPassword.addEventListener('input', function() {
        const password = this.value;
        const confirm = regConfirmPassword.value;
        
        showPasswordStrength(password);
        
        // Проверяем совпадение, если поле подтверждения не пустое
        if (confirm) {
            checkPasswordMatch(password, confirm);
        }
        
        // Обновляем стиль поля
        if (password.length > 0) {
            const validation = validatePassword(password);
            if (validation.errors.length === 0) {
                this.classList.remove('error-field');
                this.classList.add('success-field');
            } else {
                this.classList.add('error-field');
                this.classList.remove('success-field');
            }
        } else {
            this.classList.remove('error-field', 'success-field');
        }
    });
    
    regConfirmPassword.addEventListener('input', function() {
        const password = regPassword.value;
        const confirm = this.value;
        
        checkPasswordMatch(password, confirm);
    });
}

// Проверка совпадения паролей
function checkPasswordMatch(password, confirm) {
    if (!passwordMatch) return;
    
    if (confirm.length === 0) {
        passwordMatch.style.display = 'none';
        regConfirmPassword.classList.remove('error-field', 'success-field');
        return;
    }
    
    if (password === confirm) {
        passwordMatch.innerHTML = '✅ Пароли совпадают';
        passwordMatch.style.color = '#2ecc71';
        passwordMatch.style.display = 'block';
        regConfirmPassword.classList.remove('error-field');
        regConfirmPassword.classList.add('success-field');
    } else {
        passwordMatch.innerHTML = '❌ Пароли не совпадают';
        passwordMatch.style.color = '#e74c3c';
        passwordMatch.style.display = 'block';
        regConfirmPassword.classList.add('error-field');
        regConfirmPassword.classList.remove('success-field');
    }
}

// Функция показа уведомлений
function showNotification(message, type = 'info') {
    // Удаляем старое уведомление, если есть
    const oldNotification = document.querySelector('.custom-notification');
    if (oldNotification) {
        oldNotification.remove();
    }
    
    // Создаем элемент уведомления
    const notification = document.createElement('div');
    notification.className = 'custom-notification';
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 25px;
        border-radius: 8px;
        color: white;
        font-weight: bold;
        z-index: 10000;
        animation: slideIn 0.3s ease;
        box-shadow: 0 6px 12px rgba(0,0,0,0.15);
        max-width: 400px;
        word-wrap: break-word;
        backdrop-filter: blur(10px);
    `;
    
    // Цвета в зависимости от типа
    if (type === 'success') {
        notification.style.backgroundColor = 'rgba(76, 175, 80, 0.95)';
        notification.style.borderLeft = '5px solid #388E3C';
    } else if (type === 'error') {
        notification.style.backgroundColor = 'rgba(244, 67, 54, 0.95)';
        notification.style.borderLeft = '5px solid #D32F2F';
    } else if (type === 'warning') {
        notification.style.backgroundColor = 'rgba(255, 152, 0, 0.95)';
        notification.style.borderLeft = '5px solid #F57C00';
    } else {
        notification.style.backgroundColor = 'rgba(33, 150, 243, 0.95)';
        notification.style.borderLeft = '5px solid #1976D2';
    }
    
    // Добавляем стили для анимации
    const style = document.createElement('style');
    style.textContent = `
        @keyframes slideIn {
            from { transform: translateX(100%); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        @keyframes fadeOut {
            from { opacity: 1; transform: translateX(0); }
            to { opacity: 0; transform: translateX(100%); }
        }
    `;
    document.head.appendChild(style);
    
    // Добавляем уведомление на страницу
    document.body.appendChild(notification);
    
    // Автоматическое удаление через 5 секунд
    setTimeout(() => {
        notification.style.animation = 'fadeOut 0.3s ease';
        setTimeout(() => {
            if (document.body.contains(notification)) {
                document.body.removeChild(notification);
            }
            if (document.head.contains(style)) {
                document.head.removeChild(style);
            }
        }, 300);
    }, 5000);
    
    // Возможность закрыть уведомление кликом
    notification.addEventListener('click', function() {
        notification.style.animation = 'fadeOut 0.3s ease';
        setTimeout(() => {
            if (document.body.contains(notification)) {
                document.body.removeChild(notification);
            }
            if (document.head.contains(style)) {
                document.head.removeChild(style);
            }
        }, 300);
    });
}

// ОБНОВЛЕННЫЙ ВХОД ИЗ РАБОЧЕЙ ВЕРСИИ
// Обработка формы входа из хедера
if (loginForm) {
    console.log('Login form found in header');
    
    // Обработчик для кнопки "войти"
    const submitBtn = loginForm.querySelector('.submit-btn');
    if (submitBtn) {
        submitBtn.addEventListener('click', async function(event) {
            event.preventDefault(); // Предотвращаем стандартную отправку формы
            
            console.log('Login button clicked');
            
            const usernameInput = loginForm.querySelector('.login-input');
            const passwordInput = loginForm.querySelector('.password-input');
            
            if (!usernameInput || !passwordInput) {
                console.error('Input fields not found');
                showNotification('Ошибка: поля ввода не найдены', 'error');
                return;
            }
            
            const username = usernameInput.value.trim();
            const password = passwordInput.value;
            
            console.log('Attempting login with username:', username);
            
            // Базовая валидация
            if (!username) {
                showNotification('Введите имя пользователя', 'error');
                usernameInput.focus();
                return;
            }
            
            if (!password) {
                showNotification('Введите пароль', 'error');
                passwordInput.focus();
                return;
            }
            
            // Проверяем соединение с сервером
            const isConnected = await checkServerConnection();
            if (!isConnected) {
                showNotification('Сервер не отвечает. Запустите бэкенд командой: go run main.go', 'error');
                return;
            }
            
            // Показываем индикатор загрузки
            const originalText = submitBtn.value;
            submitBtn.value = 'Вход...';
            submitBtn.disabled = true;
            
            try {
                console.log('Sending login request to:', API_URL + '/login');
                console.log('Request payload:', { username, password: '***' });
                
                const response = await fetch(API_URL + '/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ 
                        username: username,
                        password: password 
                    })
                });
                
                console.log('Login response status:', response.status);
                
                const responseText = await response.text();
                console.log('Login response text:', responseText);
                
                let data;
                try {
                    data = JSON.parse(responseText);
                    console.log('Login response JSON:', data);
                } catch (parseError) {
                    console.error('Failed to parse JSON:', parseError);
                    showNotification('Ошибка сервера: неверный формат ответа', 'error');
                    return;
                }
                
                if (response.ok && data.success) {
                    showNotification('✅ Вход выполнен успешно!', 'success');
                    // Сохраняем токен и данные пользователя
                    saveAuthData(data.token, data.user);
                    if (response.ok && data.success) {
                    showNotification('✅ Вход выполнен успешно!', 'success');
                    saveAuthData(data.token, data.user);
                    
                    // ОЧИСТКА ПОЛЕЙ ВВОДА (вместо reset)
                    const usernameInput = loginForm.querySelector('.login-input');
                    const passwordInput = loginForm.querySelector('.password-input');
                    
                    if (usernameInput) usernameInput.value = '';
                    if (passwordInput) passwordInput.value = '';
                    
                    updateUIAfterLogin(data.user.username);
}
                    updateUIAfterLogin(data.user.username);
                } else {
                    const errorMessage = data.message || 'Неверный логин или пароль';
                    showNotification(`❌ ${errorMessage}`, 'error');
                }
            } catch (error) {
                console.error('Login fetch error:', error);
                console.error('Error details:', {
                    name: error.name,
                    message: error.message,
                    stack: error.stack
                });
                showNotification(`❌ Ошибка сети: ${error.message}. Проверьте, запущен ли сервер на localhost:8080`, 'error');
            } finally {
                // Восстанавливаем кнопку
                submitBtn.value = originalText;
                submitBtn.disabled = false;
            }
        });
    }
    
    // Также можно добавить обработчик нажатия Enter
    const loginInput = loginForm.querySelector('.login-input');
    const passwordInput = loginForm.querySelector('.password-input');
    
    if (loginInput && passwordInput) {
        [loginInput, passwordInput].forEach(input => {
            input.addEventListener('keypress', function(event) {
                if (event.key === 'Enter') {
                    event.preventDefault();
                    if (submitBtn) submitBtn.click();
                }
            });
        });
    }
}

// Обработка формы входа в модальном окне (если оно используется)
if (loginModalForm) {
    loginModalForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        
        const username = loginUsername.value.trim();
        const password = loginPassword.value;
        
        // Валидация
        if (!username) {
            showNotification('Введите имя пользователя', 'error');
            loginUsername.focus();
            return;
        }
        
        if (!password) {
            showNotification('Введите пароль', 'error');
            loginPassword.focus();
            return;
        }
        
        // Показываем индикатор загрузки
        const submitBtn = loginModalForm.querySelector('.modal-submit');
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '⌛ Вход...';
        submitBtn.disabled = true;
        
        try {
            const response = await fetch(`${API_URL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password })
            });
            
            console.log('Modal login response status:', response.status);
            
            const responseText = await response.text();
            let data;
            try {
                data = JSON.parse(responseText);
            } catch (parseError) {
                console.error('Failed to parse JSON:', parseError);
                showNotification('Ошибка сервера: неверный формат ответа', 'error');
                return;
            }
            
            if (response.ok && data.success) {
                showNotification('✅ Вход выполнен успешно!', 'success');
                // Сохраняем токен и данные пользователя
                saveAuthData(data.token, data.user);
                loginModal.style.display = 'none';
                loginModalForm.reset();
                updateUIAfterLogin(data.user.username);
            } else {
                const errorMessage = data.message || 'Неверный логин или пароль';
                showNotification(`❌ ${errorMessage}`, 'error');
            }
        } catch (error) {
            console.error('Modal login error:', error);
            showNotification('❌ Ошибка при входе. Проверьте подключение к серверу.', 'error');
        } finally {
            // Восстанавливаем кнопку
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    });
}

// Обработка формы регистрации
if (registerForm) {
    console.log('Registration form found');
    
    registerForm.addEventListener('submit', async function(event) {
        event.preventDefault();
        
        const username = regUsername.value.trim();
        const email = regEmail.value.trim();
        const password = regPassword.value;
        const confirmPassword = regConfirmPassword.value;
        
        console.log('Registration attempt:', { username, email, password: '***' });
        
        // Валидация полей
        if (!validateUsernameField(regUsername)) {
            showNotification('Имя пользователя должно содержать минимум 3 символа и может содержать только буквы, цифры, точки, дефисы и подчеркивания', 'error');
            regUsername.focus();
            return;
        }
        
        if (!validateEmailField(regEmail)) {
            showNotification('Введите корректный email адрес', 'error');
            regEmail.focus();
            return;
        }
        
        if (!password) {
            showNotification('Введите пароль', 'error');
            regPassword.focus();
            return;
        }
        
        if (password !== confirmPassword) {
            showNotification('Пароли не совпадают', 'error');
            regConfirmPassword.focus();
            return;
        }
        
        // Проверка сложности пароля
        const validation = validatePassword(password);
        if (validation.errors.length > 0) {
            showNotification('Пароль не соответствует требованиям безопасности. Проверьте список ошибок.', 'error');
            regPassword.focus();
            return;
        }
        
        // Дополнительная проверка на сервере
        const serverValidation = await validatePasswordOnServer(password);
        if (!serverValidation.valid) {
            showNotification('Пароль не соответствует требованиям безопасности', 'error');
            regPassword.focus();
            return;
        }
        
        // Проверяем соединение с сервером
        const isConnected = await checkServerConnection();
        if (!isConnected) {
            showNotification('Сервер не отвечает. Запустите бэкенд командой: go run main.go', 'error');
            return;
        }
        
        // Показываем индикатор загрузки
        const submitBtn = document.getElementById('submitBtn');
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '⌛ Регистрация...';
        submitBtn.disabled = true;
        
        // Отправка данных на сервер
        try {
            console.log('Sending registration request to:', API_URL + '/register');
            
            const response = await fetch(API_URL + '/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    username, 
                    email, 
                    password,
                    confirm_password: confirmPassword 
                })
            });
            
            console.log('Registration response status:', response.status);
            
            const responseText = await response.text();
            console.log('Registration response text:', responseText);
            
            let data;
            try {
                data = JSON.parse(responseText);
                console.log('Registration response JSON:', data);
            } catch (parseError) {
                console.error('Failed to parse JSON:', parseError);
                showNotification('Ошибка сервера: неверный формат ответа', 'error');
                return;
            }
            
            if (response.ok && data.success) {
                showNotification('✅ Регистрация успешно завершена!', 'success');
                // Сохраняем токен и данные пользователя
                saveAuthData(data.token, data.user);
                modal.style.display = 'none';
                registerForm.reset();
                resetPasswordValidation();
                resetFieldStyles();
                updateUIAfterLogin(data.user.username);
            } else {
                const errorMessage = data.message || data.errors?.join(', ') || 'Ошибка при регистрации';
                showNotification(`❌ ${errorMessage}`, 'error');
            }
        } catch (error) {
            console.error('Registration fetch error:', error);
            showNotification('❌ Ошибка при регистрации. Проверьте подключение к серверу.', 'error');
        } finally {
            // Восстанавливаем кнопку
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    });
}

// Добавляем стили для кнопки "Забыли пароль?"
if (document.getElementById('forgotPassword')) {
    document.getElementById('forgotPassword').addEventListener('click', function(event) {
        event.preventDefault();
        showNotification('Функция восстановления пароля временно недоступна', 'warning');
    });
}

function saveAuthData(token, user) {
    if (token) {
        localStorage.setItem('auth_token', token);
        console.log('Token saved to localStorage');
    }

    if (user) {
        if (user.username) {
            localStorage.setItem('username', user.username);
        }
        if (user.email) {
            localStorage.setItem('user_email', user.email);
        }

        // 🟢 САМОЕ ВАЖНОЕ: сохраняем ID пользователя для корзины
        if (user.id !== undefined && user.id !== null) {
            localStorage.setItem('user_id', user.id);
            localStorage.setItem('currentUserId', user.id.toString());
            console.log('User ID saved to localStorage:', user.id);
        } else {
            console.warn('user.id не пришёл с сервера, currentUserId не установлен');
        }

        console.log('User data saved to localStorage:', user);
    }
}

// Проверка, авторизован ли пользователь
function isUserLoggedIn() {
    const token = localStorage.getItem('auth_token');
    const username = localStorage.getItem('username');
    // при желании можно ещё проверять currentUserId / user_id
    return !!(token && username);
}


// Функция обновления интерфейса после входа
function updateUIAfterLogin(username) {
    const loginSection = document.querySelector('.login-section');
    if (!loginSection) {
        console.error('Login section not found');
        return;
    }
    
    console.log('Updating UI for user:', username);
    
    // Проверяем, есть ли уже форма приветствия
    if (loginSection.querySelector('.welcome-section')) {
        console.log('Welcome section already exists');
        return;
    }
    
    // Создаем блок приветствия
    const welcomeDiv = document.createElement('div');
    welcomeDiv.className = 'welcome-section';
    welcomeDiv.innerHTML = `
        <div style="display: flex; flex-direction: column; align-items: flex-end;">
            <p style="margin: 0 0 5px 0; color: white; font-size: 14px;">
                👋 Привет, <strong style="color: #ffd700;">${username}</strong>!
            </p>
            <button id="logoutBtn" style="
                background: linear-gradient(135deg, #f44336, #d32f2f);
                color: white;
                border: none;
                padding: 6px 16px;
                border-radius: 4px;
                cursor: pointer;
                font-size: 13px;
                font-weight: bold;
                transition: all 0.3s;
                box-shadow: 0 2px 4px rgba(0,0,0,0.2);
                margin-top: 5px;
            ">
                🚪 Выйти
            </button>
        </div>
    `;
    
    // Заменяем форму входа
    loginSection.innerHTML = '';
    loginSection.appendChild(welcomeDiv);
    
    // Добавляем обработчик для кнопки выхода
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('mouseover', function() {
            this.style.transform = 'translateY(-2px)';
            this.style.boxShadow = '0 4px 8px rgba(0,0,0,0.3)';
        });
        
        logoutBtn.addEventListener('mouseout', function() {
            this.style.transform = 'translateY(0)';
            this.style.boxShadow = '0 2px 4px rgba(0,0,0,0.2)';
        });
        
        logoutBtn.addEventListener('click', logout);
        console.log('Logout button added');
    }
}

function logout() {
    showNotification('👋 До свидания! Вы вышли из системы.', 'info');

    // чистим данные авторизации
    localStorage.removeItem('auth_token');
    localStorage.removeItem('username');
    localStorage.removeItem('user_email');

    // 🟢 ВАЖНО: отвязываем текущего пользователя от корзины
    localStorage.removeItem('currentUserId');

    // ⛔ НЕ трогаем cart_* — корзины остаются, чтобы при следующем входе под тем же пользователем они подтянулись

    console.log('User logged out, auth data cleared, currentUserId removed');
    setTimeout(() => {
        location.reload();
    }, 1500);
}

// Ограничение формы отзыва только для авторизованных пользователей
function initFeedbackAccess() {
    const feedbackForm = document.querySelector('.feedback-form');
    if (!feedbackForm) return; // мы не на странице отзыва

    const warningId = 'feedback-auth-warning';

    function showAuthWarning() {
        // если уже есть предупреждение — не дублируем
        if (document.getElementById(warningId)) return;

        const warning = document.createElement('div');
        warning.id = warningId;
        warning.style.marginTop = '15px';
        warning.style.padding = '10px 15px';
        warning.style.borderRadius = '6px';
        warning.style.backgroundColor = '#fff3cd';
        warning.style.border = '1px solid #ffeeba';
        warning.style.color = '#856404';
        warning.style.fontSize = '14px';

        warning.innerHTML = `
            🔒 Оставлять отзывы могут только авторизованные пользователи.<br>
            Пожалуйста, войдите или зарегистрируйтесь через форму вверху страницы.
        `;

        feedbackForm.parentNode.insertBefore(warning, feedbackForm);
    }

    // при загрузке страницы сразу проверим
    if (!isUserLoggedIn()) {
        showAuthWarning();
    }

    // перехватываем отправку формы
    feedbackForm.addEventListener('submit', function (event) {
        if (!isUserLoggedIn()) {
            event.preventDefault();
            showNotification(
                'Только авторизованные пользователи могут оставлять отзывы. Войдите или зарегистрируйтесь.',
                'warning'
            );

            // чуть прокрутим к шапке, чтобы юзер увидел форму логина
            window.scrollTo({ top: 0, behavior: 'smooth' });
        } else {
            // здесь можно позже добавить отправку отзыва на сервер
            // сейчас просто даём форме работать, как задумано
        }
    });
}


document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, checking auth...');
    
    const username = localStorage.getItem('username');
    const token = localStorage.getItem('auth_token');
    
    console.log('Stored username:', username);
    console.log('Stored token:', token ? 'exists' : 'not found');
    
    if (username && token) {
        console.log('User is logged in, updating UI...');
        updateUIAfterLogin(username);
    } else {
        console.log('User is not logged in');
    }
    
    // 🔹 Ограничиваем доступ к форме отзыва
    initFeedbackAccess();

    // Проверка соединения с сервером
    checkServerConnection();
    
    // Загружаем правила паролей
    loadPasswordRules();
});


// Добавляем стили для отображения правил пароля
const passwordRulesStyles = document.createElement('style');
passwordRulesStyles.textContent = `
    .password-requirements {
        background: #f8f9fa;
        border: 1px solid #dee2e6;
        border-radius: 5px;
        padding: 15px;
        margin: 15px 0;
        font-size: 14px;
    }
    
    .password-requirements h4 {
        margin-top: 0;
        color: #333;
        font-size: 16px;
    }
    
    .password-requirements ul {
        margin: 10px 0;
        padding-left: 20px;
    }
    
    .password-requirements li {
        margin-bottom: 5px;
        color: #666;
    }
    
    .password-requirements li.valid {
        color: #28a745;
    }
    
    .password-requirements li.invalid {
        color: #dc3545;
    }
    
    .password-strength-meter {
        height: 10px;
        background: #e9ecef;
        border-radius: 5px;
        margin: 10px 0;
        overflow: hidden;
    }
    
    .password-strength-meter-fill {
        height: 100%;
        transition: width 0.3s;
        border-radius: 5px;
    }
    
    .strength-weak { background: #dc3545; }
    .strength-fair { background: #ffc107; }
    .strength-good { background: #28a745; }
    .strength-strong { background: #20c997; }
    .strength-very-strong { background: #007bff; }
`;
document.head.appendChild(passwordRulesStyles);

// Добавим кнопку для тестирования соединения
const testBtn = document.createElement('button');
testBtn.textContent = '🔧 Тест соединения';
testBtn.style.cssText = `
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 10px 15px;
    background: #666;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    z-index: 9999;
    font-size: 12px;
`;
testBtn.addEventListener('click', checkServerConnection);
document.body.appendChild(testBtn);

async function loadCarData() {
    try {
        const response = await fetch('/api/cars');
        const cars = await response.json();
        
        // Очищаем существующие карточки
        const sections = document.querySelectorAll('.cards-grid');
        sections.forEach(section => section.innerHTML = '');
        
        // Создаем карточки для каждого автомобиля
        cars.forEach(car => {
            createCarCard(car);
        });
    } catch (error) {
        console.error('Error loading car data:', error);
        // В случае ошибки показываем статические данные
        console.log('Using static data');
    }
}

// Функция создания карточки автомобиля
function createCarCard(car) {
    // Находим нужную секцию по категории
    const section = document.getElementById(car.category);
    if (!section) return;
    
    const cardsGrid = section.querySelector('.cards-grid');
    if (!cardsGrid) return;
    
    const card = document.createElement('div');
    card.className = 'card';
    
    card.innerHTML = `
        <img class="card-image" src="${car.image}" alt="${car.title}">
        <div class="card-content">
            <h3 class="card-title">${car.title}</h3>
            <div class="card-price">${car.price}</div>
            <button class="card-button" onclick="showCarDetails('${car.model}')">Подробнее</button>
        </div>
    `;
    
    cardsGrid.appendChild(card);
}

// Функция показа деталей автомобиля (обновленная)
async function showCarDetails(carModel) {
    try {
        const response = await fetch(`/api/car?model=${carModel}`);
        const data = await response.json();
        
        const car = data.car;
        const details = data.details;
        
        // Заполняем модальное окно данными
        document.getElementById('modalCarTitle').textContent = car.title;
        document.getElementById('modalPrice').textContent = car.price;
        document.getElementById('modalImage').src = car.image;
        document.getElementById('modalImage').alt = car.title;
        document.getElementById('modalDescription').textContent = car.description;

        // Заполняем технические характеристики
        const techSpecsList = document.getElementById('techSpecs');
        techSpecsList.innerHTML = '';
        details.techSpecs.forEach(spec => {
            const li = document.createElement('li');
            li.innerHTML = `<span class="spec-name">${spec.name}</span><span class="spec-value">${spec.value}</span>`;
            techSpecsList.appendChild(li);
        });

        // Заполняем комплектацию
        const equipmentList = document.getElementById('equipmentSpecs');
        equipmentList.innerHTML = '';
        details.equipment.forEach(item => {
            const li = document.createElement('li');
            li.innerHTML = `<span class="spec-name">${item.name}</span><span class="spec-value">${item.value}</span>`;
            equipmentList.appendChild(li);
        });

        // Заполняем особенности
        const featuresGrid = document.getElementById('featuresGrid');
        featuresGrid.innerHTML = '';
        details.features.forEach(feature => {
            const div = document.createElement('div');
            div.className = 'feature-item';
            div.textContent = feature;
            featuresGrid.appendChild(div);
        });

        // Показываем модальное окно
        modal.style.display = 'block';
        
    } catch (error) {
        console.error('Error loading car details:', error);
        alert('Ошибка при загрузке данных автомобиля');
    }
}

// Загружаем данные при загрузке страницы
document.addEventListener('DOMContentLoaded', loadCarData);

function openPrivacyModal() {
            document.getElementById('privacyModal').style.display = 'block';
        }
        
        function openTermsModal() {
            document.getElementById('termsModal').style.display = 'block';
        }
        
        // Закрытие модальных окон
        document.querySelectorAll('.modal-close').forEach(closeBtn => {
            closeBtn.addEventListener('click', function() {
                this.closest('.modal').style.display = 'none';
            });
        });
        
        // Закрытие при клике вне окна
        window.addEventListener('click', function(event) {
            if (event.target.classList.contains('modal')) {
                event.target.style.display = 'none';
            }
        });
        
        // Закрытие на Escape
        document.addEventListener('keydown', function(event) {
            if (event.key === 'Escape') {
                document.querySelectorAll('.modal').forEach(modal => {
                    modal.style.display = 'none';
                });
            }
        });