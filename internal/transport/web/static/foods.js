async function jfetch(url, opts={}) {
    const res = await fetch(url, {headers:{'Content-Type':'application/json'}, ...opts});
    if (!res.ok) throw new Error(await res.text());
    return res.status === 204 ? null : res.json();
  }
  
  const tbody = document.getElementById('foods-tbody');
  const catSelect = document.getElementById('food-cat-select');
  let foods = [], cats = [];
  
  async function load() {
    [foods, cats] = await Promise.all([jfetch('/api/foods'), jfetch('/api/categories')]);
    catSelect.innerHTML = '<option value="">Категория…</option>' +
      cats.map(c=>`<option value="${c.id}">${c.name}</option>`).join('');
    draw();
  }
  
  function draw() {
    tbody.innerHTML = foods.length ? foods.map(f=>`
      <tr>
        <td>${f.id}</td>
        <td contenteditable data-field="name" data-id="${f.id}">${f.name}</td>
        <td>
          <select data-field="category_id" data-id="${f.id}">
            ${cats.map(c=>`<option value="${c.id}" ${c.id===f.category_id?'selected':''}>${c.name}</option>`).join('')}
          </select>
        </td>
        <td><input type="checkbox" data-field="is_complex" data-id="${f.id}" ${f.is_complex?'checked':''}></td>
        <td>
          <button data-action="save" data-id="${f.id}">Сохранить</button>
          <button data-action="delete" data-id="${f.id}">Удалить</button>
        </td>
      </tr>`).join('')
    : `<tr><td colspan="5" class="muted" style="text-align:center">Пока пусто</td></tr>`;
  }
  
  document.getElementById('food-form').addEventListener('submit', async e=>{
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const payload = {name: fd.get('name'), category_id:+fd.get('category_id'), is_complex:!!fd.get('is_complex')};
    const created = await jfetch('/api/foods', {method:'POST', body:JSON.stringify(payload)});
    foods.push(created); e.currentTarget.reset(); draw();
  });
  
  tbody.addEventListener('click', async e=>{
    const id = +e.target.dataset.id;
    if (e.target.dataset.action==='delete') {
      await jfetch(`/api/foods/${id}`, {method:'DELETE'});
      foods = foods.filter(f=>f.id!==id); draw();
    } else if (e.target.dataset.action==='save') {
      const tr = e.target.closest('tr');
      const name = tr.querySelector('[data-field="name"]').textContent.trim();
      const cat  = +tr.querySelector('[data-field="category_id"]').value;
      const isC  = tr.querySelector('[data-field="is_complex"]').checked;
      const upd = await jfetch(`/api/foods/${id}`, {method:'PUT', body:JSON.stringify({name,category_id:cat,is_complex:isC})});
      foods[foods.findIndex(f=>f.id===id)] = upd; draw();
    }
  });
  
  load();
  