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

(function () {
    if (window.__appInit2) return; window.__appInit2 = true;
    function $(id) { return document.getElementById(id); }

    function splitCmds(val) {
        return val.split(/\r?\n/).map(s => s.trim()).filter(Boolean);
    }

    // --- run-cmds ---
    async function onRunCmds(e) {
        e.preventDefault();
        const lab = $('lab').value.trim();
        const sudo = $('sudo').checked;
        const timeoutSec = parseInt($('timeout').value, 10) || 15;
        const user = $('euser').value;
        const pass = $('epass').value;
        const format = $('fmt').value;
        const cmds = splitCmds($('cmds').value);

        const res = await fetch('/run-cmds', {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ lab, sudo, timeoutSec, user, pass, format, cmds })
        });
        const data = await res.json().catch(() => ({ ok: false, error: 'bad json' }));
        if (!res.ok || !data.ok) {
            alert('Run error: ' + (data.error || res.statusText)); return;
        }

        const wrap = $('execOut'); const div = $('execResults');
        div.innerHTML = '';
        (data.results || []).forEach(r => {
            const pre = document.createElement('pre');
            let bodyPretty = '';
            try { bodyPretty = JSON.stringify(JSON.parse(r.body || '{}'), null, 2) } catch { bodyPretty = (r.body || '') + ''; }
            pre.textContent = [
                `${r.name} (${r.ip}) [${r.kind}]`,
                `OK=${r.ok} status=${r.status}${r.error ? ' error=' + r.error : ''}`,
                bodyPretty
            ].join('\n');
            div.appendChild(pre);
        });
        wrap.hidden = false;
    }

    // --- health ---
    async function onHealth(e) {
        e.preventDefault();
        const lab = $('lab').value.trim();
        const sudo = $('sudo').checked;
        const timeoutSec = parseInt($('timeout').value, 10) || 20;
        const user = $('euser').value;
        const pass = $('epass').value;

        const res = await fetch('/health', {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ lab, sudo, timeoutSec, user, pass })
        });
        const data = await res.json().catch(() => ({ ok: false, error: 'bad json' }));
        if (!res.ok || !data.ok) {
            alert('Health error: ' + (data.error || res.statusText)); return;
        }

        const out = $('healthOut'); const div = $('healthResults');
        div.innerHTML = '';
        (data.nodes || []).forEach(n => {
            const s = document.createElement('div');
            const lines = n.checks.map(c => `- ${c.name}: ${c.result}${c.detail ? ' (' + c.detail + ')' : ''}`).join('\n');
            s.innerHTML = `<h4>${n.name} (${n.ip})</h4><pre>${lines}</pre>`;
            div.appendChild(s);
        });
        out.hidden = false;
    }

    document.addEventListener('DOMContentLoaded', () => {
        $('execForm').addEventListener('submit', onRunCmds);
        $('healthForm').addEventListener('submit', onHealth);
    });
})();