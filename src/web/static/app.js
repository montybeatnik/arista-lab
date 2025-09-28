(function () {
    if (window.__appInit) return;      // ← guard
    window.__appInit = true;

    function $(id) { return document.getElementById(id); }

    async function onSubmit(e) {
        e.preventDefault();
        const btn = $('runBtn'), spin = $('spinner'), out = $('result');
        btn.disabled = true; spin.hidden = false; out.hidden = true;

        try {
            const body = {
                lab: $('lab').value.trim(),
                timeoutSec: parseInt($('timeout').value, 10) || 15,
                sudo: $('sudo').checked
            };
            const res = await fetch('/inspect', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body)
            });
            const data = await res.json().catch(() => ({}));
            if (!res.ok || !data.ok) { alert('Error: ' + (data.error || res.statusText)); return; }

            $('labKey').textContent = data.labKey || '(unknown)';
            const tbody = document.querySelector('#nodesTbl tbody');
            tbody.innerHTML = '';
            (data.nodes || []).forEach(n => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
            <td>${n.name}</td><td>${n.kind}</td><td>${n.image}</td>
            <td>${n.state} <span class="muted">${n.status || ''}</span></td>
            <td>${n.ipv4_address || n.ipv4 || ''}</td><td>${n.owner || ''}</td>`;
                tbody.appendChild(tr);
            });
            $('rawJson').textContent = JSON.stringify(data.rawJson ?? {}, null, 2);
            out.hidden = false;
        } catch (err) {
            alert('Request error: ' + err);
        } finally {
            btn.disabled = false; spin.hidden = true;
        }
    }

    document.addEventListener('DOMContentLoaded', () => {
        $('f').addEventListener('submit', onSubmit);
    });
})();