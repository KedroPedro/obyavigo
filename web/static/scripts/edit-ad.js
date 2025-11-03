document.addEventListener('DOMContentLoaded', function() {
  loadAdData();
  initFormSubmit();
  initImageUpload();
});

let currentImages = [];

async function loadAdData() {
  const form = document.getElementById('editAdForm');
  const adId = form.dataset.adId;
  
  if (!adId) {
    alert('Ошибка: ID объявления не найден');
    window.location.href = '/profile/';
    return;
  }

  try {
    const response = await fetch(`/api/ads/${adId}/`, {
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to load ad data');
    }

    const result = await response.json();
    
    if (result.success && result.data) {
      const ad = result.data;
      
      
      document.getElementById('adTitle').value = ad.title || '';
      document.getElementById('adPrice').value = ad.price || '';
      document.getElementById('adCondition').value = ad.condition || '';
      document.getElementById('adDescription').value = ad.desc || '';
      document.getElementById('adPhone').value = ad.contactPhone || '';
      
      
      if (ad.images && ad.images.length > 0) {
        currentImages = ad.images;
        displayCurrentImages();
      }
    } else {
      throw new Error('Ad data not found');
    }
  } catch (error) {
    console.error('Error loading ad:', error);
    alert('Ошибка при загрузке объявления');
    window.location.href = '/profile/';
  }
}

function initFormSubmit() {
  const form = document.getElementById('editAdForm');
  
  form.addEventListener('submit', async function(e) {
    e.preventDefault();
    
    clearAllErrors();
    
    if (!validateForm()) {
      return;
    }

    const adId = form.dataset.adId;
    const submitBtn = form.querySelector('.submit-btn');
    const originalText = submitBtn.textContent;
    
    submitBtn.disabled = true;
    submitBtn.textContent = 'Сохранение...';

    const data = {
      title: document.getElementById('adTitle').value.trim(),
      description: document.getElementById('adDescription').value.trim(),
      price: parseInt(document.getElementById('adPrice').value),
      condition: document.getElementById('adCondition').value,
      contact_phone: document.getElementById('adPhone').value.trim() || '',
    };

    try {
      const response = await fetch(`/api/ads/${adId}/`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      const result = await response.json();

      if (response.ok) {
        alert('Объявление успешно обновлено!');
        window.location.href = '/profile/';
      } else {
        alert(result.message || 'Ошибка при обновлении объявления');
        submitBtn.disabled = false;
        submitBtn.textContent = originalText;
      }
    } catch (error) {
      console.error('Error updating ad:', error);
      alert('Ошибка подключения к серверу');
      submitBtn.disabled = false;
      submitBtn.textContent = originalText;
    }
  });
}

function validateForm() {
  let isValid = true;

  const title = document.getElementById('adTitle').value.trim();
  if (!title) {
    showError('titleError', 'Введите заголовок');
    isValid = false;
  } else if (title.length < 5) {
    showError('titleError', 'Заголовок должен быть не менее 5 символов');
    isValid = false;
  }

  const price = document.getElementById('adPrice').value;
  if (!price || price <= 0) {
    showError('priceError', 'Укажите корректную цену');
    isValid = false;
  }

  const description = document.getElementById('adDescription').value.trim();
  if (!description) {
    showError('descriptionError', 'Введите описание');
    isValid = false;
  } else if (description.length < 20) {
    showError('descriptionError', 'Описание должно быть не менее 20 символов');
    isValid = false;
  }

  const condition = document.getElementById('adCondition').value;
  if (!condition) {
    showError('conditionError', 'Выберите состояние товара');
    isValid = false;
  }

  const phone = document.getElementById('adPhone').value.trim();
  if (phone && !isValidBelarusPhone(phone)) {
    showError('phoneError', 'Неверный формат белорусского номера');
    isValid = false;
  }

  return isValid;
}

function showError(elementId, message) {
  const errorEl = document.getElementById(elementId);
  if (errorEl) {
    errorEl.textContent = message;
  }
}

function clearAllErrors() {
  const errors = document.querySelectorAll('.error-message');
  errors.forEach(el => el.textContent = '');
}

function isValidBelarusPhone(phone) {
  if (!phone) return true;
  
  const cleaned = phone.replace(/\D/g, '');
  
  if (cleaned.length === 12 && cleaned.startsWith('375')) {
    return /^375(25|29|33|44)\d{7}$/.test(cleaned);
  }
  
  if (cleaned.length === 9) {
    return /^(25|29|33|44)\d{7}$/.test(cleaned);
  }
  
  return false;
}

function displayCurrentImages() {
  const container = document.getElementById('currentImages');
  container.innerHTML = '';
  
  currentImages.forEach((imageId) => {
    const imageItem = document.createElement('div');
    imageItem.className = 'image-item';
    imageItem.innerHTML = `
      <img src="/api/images/${imageId}/" alt="Ad image" />
      <button type="button" class="delete-btn" onclick="deleteImage('${imageId}')">×</button>
    `;
    container.appendChild(imageItem);
  });
}

async function deleteImage(imageId) {
  const form = document.getElementById('editAdForm');
  const adId = form.dataset.adId;
  
  if (!confirm('Удалить это изображение?')) {
    return;
  }
  
  try {
    const response = await fetch(`/api/ads/${adId}/images/${imageId}/`, {
      method: 'DELETE',
      credentials: 'include',
    });
    
    if (response.ok) {
      currentImages = currentImages.filter(id => id !== imageId);
      displayCurrentImages();
    } else {
      alert('Ошибка при удалении изображения');
    }
  } catch (error) {
    console.error('Error deleting image:', error);
    alert('Ошибка подключения к серверу');
  }
}

function initImageUpload() {
  const imageInput = document.getElementById('imageInput');
  
  imageInput.addEventListener('change', async function(e) {
    const files = e.target.files;
    if (files.length === 0) return;
    
    const form = document.getElementById('editAdForm');
    const adId = form.dataset.adId;
    
    const formData = new FormData();
    for (let file of files) {
      formData.append('images', file);
    }
    
    try {
      const response = await fetch(`/api/ads/${adId}/images/`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      
      const result = await response.json();
      
      if (response.ok && result.image_ids) {
        currentImages = currentImages.concat(result.image_ids);
        displayCurrentImages();
        imageInput.value = '';
      } else {
        alert(result.message || 'Ошибка при загрузке изображений');
      }
    } catch (error) {
      console.error('Error uploading images:', error);
      alert('Ошибка подключения к серверу');
    }
  });
}
