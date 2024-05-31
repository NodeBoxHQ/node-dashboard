function attachTippy() {
    if(document.getElementById('activity-bar')) {
        tippy.setDefaultProps({ maxWidth: '30em' })
        tippy('#activity-bar', {
            popperOptions: {
                modifiers: [
                    {
                        name: 'flip',
                        options: {
                            fallbackPlacements: ['top'],
                        },
                    },
                ],
            },
            placement: (nodeType === "Xally") ? 'top' : undefined,
            flip: false,
            content: document.getElementById('activity-bar').getAttribute('tooltip-data'),
            allowHTML: true,
            interactive: true,
        });
    }
}

document.addEventListener('htmx:load', function(evt) {
    attachTippy();
});

document.addEventListener("htmx:confirm", function(e) {
    e.preventDefault();
    if (!e.target.hasAttribute('hx-confirm')) {
        e.detail.issueRequest(true);
        return;
    }

    console.log(e.detail.question.includes("node"), nodeType, e.detail.question)

    if (e.detail.question.includes("node") && (nodeType === "Xally" || nodeType === "Juneo")) {
        Swal.fire({
            title: "Xally does not support restarting",
            icon: "error",
            toast: true,
            position: "bottom",
            showConfirmButton: false,
            timer: 3000,
            width: 400,
            timerProgressBar: true,
        });

        return;
    }


    Swal.fire({
        title: "Are you sure?",
        text: `${e.detail.question}`,
        icon: "warning",
        showCancelButton: true,
        confirmButtonText: "Yes",
        cancelButtonText: "No",
        reverseButtons: false,
        confirmButtonColor: "#dc3545",
        cancelButtonColor: "#6c757d",
    }).then(async function(result) {
        if(result.isConfirmed) {
            e.detail.issueRequest(true);
            const text = "This will take some time, please wait...";
            let title = "";

            if (e.detail.question.includes("This will restart the server, you may lose connection for some time")) {
                title = "Restarting Server";
            } else if (e.detail.question.includes("This will restart the blockchain node running on your server")) {
                title = "Restarting Blockchain Node";
            }

            await Swal.fire({
                title: title,
                text: "This will take some time, please wait...",
                icon: "info",
                showConfirmButton: false,
                allowOutsideClick: false,
                timer: 60000,
                timerProgressBar: true,
            });

            window.location.reload();
        }
    });
});