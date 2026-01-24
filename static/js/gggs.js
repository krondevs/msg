//jsonRequest(jsonDataSerialized, route, callback)
//jsonRequestNoDim(jsonDataSerialized, route, callback)
//sendMultiPartData(method, route, formId, callback)
//sendForm(method, route, formId, callback)
//formToJSON(formId)
// showNotification(message, type = "success", duration = 5000)

function comprobantes(){
  let datos = {
    "sesion":sessionHash
  }
  let jsonDataSerialized = JSON.stringify(datos);
  jsonRequest(jsonDataSerialized, "/comprobantes", (data)=>{
    let text = "";
    console.log(data)
    for(let i of data.data){
      text+=`<option value="${i["key"]}">${i["value"]}</option>`
    }
    $("#main").html(`
      <h4>Seleccione el comprobante</h4>
      <form id="forma" method="POST" enctype="multipart/form-data">
      <select id="comprobante" onchange="tipoFormulario()" class="form-control" name="comprobante">
      ${text}
      </select>
      <input type="hidden" name="sesion" value="${sessionHash}">
      <br><br>
      <button class="btn btn-success" type="button" onclick="tipoFormulario()">Cargar Formulario</button>
      </form>
      <hr>
      <div id="carga"></div>
      `);
  });
}

function tipoFormulario(){
  let plantilla = $("#comprobante").val();
  let text = "";
  let e = plantilla;
  switch (e){
    case "1.html":
      text = `
       <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Remittance Transfer Receipt Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="AcountNumber" class="form-label">Account Number</label>
                                    <input type="text" class="form-control" id="AcountNumber" name="AcountNumber">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="TransferDate" class="form-label">Transfer Date</label>
                                    <input type="text" class="form-control" id="TransferDate" name="TransferDate">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Sender" class="form-label">Sender</label>
                                    <input type="text" class="form-control" id="Sender" name="Sender">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Reference" class="form-label">Reference</label>
                                    <input type="text" class="form-control" id="Reference" name="Reference">
                                </div>
                            </div>

                            <div class="mb-3">
                                <label for="SenderAddress" class="form-label">Sender Address</label>
                                <textarea class="form-control" id="SenderAddress" name="SenderAddress" rows="2"></textarea>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankName" class="form-label">Bank Name</label>
                                    <input type="text" class="form-control" id="BankName" name="BankName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="AccountName" class="form-label">Account Name</label>
                                    <input type="text" class="form-control" id="AccountName" name="AccountName">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankAccount" class="form-label">Bank Account (last four)</label>
                                    <input type="text" class="form-control" id="BankAccount" name="BankAccount">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="RoutingId" class="form-label">Routing ID</label>
                                    <input type="text" class="form-control" id="RoutingId" name="RoutingId">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="TransferAmount" class="form-label">Transfer Amount</label>
                                    <input type="text" class="form-control" id="TransferAmount" name="TransferAmount">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="TransferFees" class="form-label">Transfer Fees</label>
                                    <input type="text" class="form-control" id="TransferFees" name="TransferFees">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Total" class="form-label">Total</label>
                                    <input type="text" class="form-control" id="Total" name="Total">
                                </div>
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                                <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
      `
    break;
    case "2.html":
      text = `
      <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Remittance Transfer Receipt Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="AcountNumber" class="form-label">Account Number</label>
                                    <input type="text" class="form-control" id="AcountNumber" name="AcountNumber">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="TransferDate" class="form-label">Transfer Date</label>
                                    <input type="text" class="form-control" id="TransferDate" name="TransferDate">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Sender" class="form-label">Sender</label>
                                    <input type="text" class="form-control" id="Sender" name="Sender">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Reference" class="form-label">Reference</label>
                                    <input type="text" class="form-control" id="Reference" name="Reference">
                                </div>
                            </div>

                            <div class="mb-3">
                                <label for="SenderAddress" class="form-label">Sender Address</label>
                                <textarea class="form-control" id="SenderAddress" name="SenderAddress" rows="2"></textarea>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Imad" class="form-label">IMAD</label>
                                    <input type="text" class="form-control" id="Imad" name="Imad">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Omad" class="form-label">OMAD</label>
                                    <input type="text" class="form-control" id="Omad" name="Omad">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankName" class="form-label">Bank Name</label>
                                    <input type="text" class="form-control" id="BankName" name="BankName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="AccountName" class="form-label">Account Name</label>
                                    <input type="text" class="form-control" id="AccountName" name="AccountName">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankAccount" class="form-label">Bank Account (last four)</label>
                                    <input type="text" class="form-control" id="BankAccount" name="BankAccount">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="RoutingId" class="form-label">Routing ID</label>
                                    <input type="text" class="form-control" id="RoutingId" name="RoutingId">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="TransferAmount" class="form-label">Transfer Amount</label>
                                    <input type="text" class="form-control" id="TransferAmount" name="TransferAmount">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="TransferFees" class="form-label">Transfer Fees</label>
                                    <input type="text" class="form-control" id="TransferFees" name="TransferFees">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Total" class="form-label">Total</label>
                                    <input type="text" class="form-control" id="Total" name="Total">
                                </div>
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                                <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
      `
      break;
      case "3.html":
        text = `
        <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-6 col-md-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Payment Details Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                         <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="mb-3">
                                <label for="MontoArriba" class="form-label">Monto</label>
                                <input type="text" class="form-control" id="MontoArriba" name="MontoArriba">
                            </div>

                            <div class="mb-3">
                                <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                            </div>

                            <div class="mb-3">
                                <label for="RoutingNumber" class="form-label">Routing Number / ABA</label>
                                <input type="text" class="form-control" id="RoutingNumber" name="RoutingNumber">
                            </div>

                            <div class="mb-3">
                                <label for="NumeroCuenta" class="form-label">Número de Cuenta</label>
                                <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                            </div>

                            <div class="mb-3">
                                <label for="Fecha" class="form-label">Fecha</label>
                                <input type="text" class="form-control" id="Fecha" name="Fecha">
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                               <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
        `
        break;
        case "4.html":
            text = `
            <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-6 col-md-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Payment Details Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                         <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="mb-3">
                                <label for="MontoArriba" class="form-label">Monto</label>
                                <input type="text" class="form-control" id="MontoArriba" name="MontoArriba">
                            </div>

                            <div class="mb-3">
                                <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                            </div>

                            <div class="mb-3">
                                <label for="RoutingNumber" class="form-label">SWIFT</label>
                                <input type="text" class="form-control" id="RoutingNumber" name="Swift">
                            </div>

                            <div class="mb-3">
                                <label for="NumeroCuenta" class="form-label">Número de Cuenta</label>
                                <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                            </div>

                            <div class="mb-3">
                                <label for="Fecha" class="form-label">Fecha</label>
                                <input type="text" class="form-control" id="Fecha" name="Fecha">
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                               <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
            `
        break;
        case "5.html":
            text = `
            <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">FEDWIRE Funds Transfer Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                         <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="Date" class="form-label">Date</label>
                                    <input type="text" class="form-control" id="Date" name="Date">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Time" class="form-label">Time</label>
                                    <input type="text" class="form-control" id="Time" name="Time">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="FedReference" class="form-label">Fed Reference #</label>
                                    <input type="text" class="form-control" id="FedReference" name="FedReference">
                                </div>
                            </div>

                            <h5 class="mt-4 mb-3">Originator Information</h5>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="OriginatorName" class="form-label">Originator Name</label>
                                    <input type="text" class="form-control" id="OriginatorName" name="OriginatorName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="OriginatorAccount" class="form-label">Originator Account</label>
                                    <input type="text" class="form-control" id="OriginatorAccount" name="OriginatorAccount">
                                </div>
                            </div>

                            <h5 class="mt-4 mb-3">Beneficiary Information</h5>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryBank" class="form-label">Beneficiary Bank</label>
                                    <input type="text" class="form-control" id="BeneficiaryBank" name="BeneficiaryBank">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryAba" class="form-label">Beneficiary ABA</label>
                                    <input type="text" class="form-control" id="BeneficiaryAba" name="BeneficiaryAba">
                                </div>
                            </div>

                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryName" class="form-label">Beneficiary Name</label>
                                    <input type="text" class="form-control" id="BeneficiaryName" name="BeneficiaryName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryAccount" class="form-label">Beneficiary Account</label>
                                    <input type="text" class="form-control" id="BeneficiaryAccount" name="BeneficiaryAccount">
                                </div>
                            </div>

                            <div class="mb-3">
                                <label for="Amount" class="form-label">Amount</label>
                                <input type="text" class="form-control" id="Amount" name="Amount">
                            </div>

                            <div class="mb-3">
                                <label for="OriginatorToBenificiaryInformation" class="form-label">Originator to Beneficiary Information (OBI)</label>
                                <textarea class="form-control" id="OriginatorToBenificiaryInformation" name="OriginatorToBenificiaryInformation" rows="3"></textarea>
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                                 <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
            `
        break;
        case "6.html":
            text = `
            <div class="container mt-4">
    <div class="row justify-content-center">
        <div class="col-lg-8">
            <div class="card">
                <div class="card-header">
                    <h4 class="mb-0">Transfer Confirmation Form</h4>
                </div>
                <div class="card-body">
                    <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">

                        <h5 class="mt-4 mb-3">Transfer Information</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Status" class="form-label">Status</label>
                                <input type="text" class="form-control" id="Status" name="Status">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="Amount" class="form-label">Amount Sent</label>
                                <input type="text" class="form-control" id="Amount" name="Amount">
                            </div>

                            <div class="col-md-6 mb-3">
                                <label for="Reference" class="form-label">Reference</label>
                                <input type="text" class="form-control" id="Reference" name="Reference">
                            </div>

                            <div class="col-md-6 mb-3">
                                <label for="Date" class="form-label">Date</label>
                                <input type="text" class="form-control" id="Date" name="Date">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="MessageToPayee" class="form-label">Message to Payee</label>
                                <input type="text" class="form-control" id="MessageToPayee" name="MessageToPayee">
                            </div>
                        </div>

                        <h5 class="mt-4 mb-3">Sender Details</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="PayeeName" class="form-label">Sender Name</label>
                                <input type="text" class="form-control" id="PayeeName" name="PayeeName">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="SenderName" class="form-label">Sent By</label>
                                <input type="text" class="form-control" id="SenderName" name="SenderName">
                            </div>
                        </div>

                        <h5 class="mt-4 mb-3">Payee Details</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="BeneficiaryName" class="form-label">Payee Name</label>
                                <input type="text" class="form-control" id="BeneficiaryName" name="BeneficiaryName">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="Account" class="form-label">Account Number</label>
                                <input type="text" class="form-control" id="Account" name="Account">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Bank" class="form-label">Bank</label>
                                <input type="text" class="form-control" id="Bank" name="Bank">
                            </div>
                             <div class="col-md-6 mb-3">
                                <label for="BankAddress" class="form-label">Bank Address</label>
                                <input type="text" class="form-control" id="BankAddress" name="BankAddress">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="BankCountry" class="form-label">Bank Country</label>
                                <input type="text" class="form-control" id="BankCountry" name="BankCountry">
                            </div>
                        </div>

                        <div class="mb-3">
                            <label for="Swift" class="form-label">SWIFT Code</label>
                            <input type="text" class="form-control" id="Swift" name="Swift">
                        </div>

                        <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                            <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
            `
        break;
        case "7.html":
            text = `
            <div class="container mt-4">
    <div class="row justify-content-center">
        <div class="col-lg-8">
            <div class="card">
                <div class="card-header">
                    <h4 class="mb-0">Instrucciones para Transferencia de Moneda Extranjera</h4>
                </div>
                <div class="card-body">
                    <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Fecha" class="form-label">Fecha</label>
                                <input type="text" class="form-control" id="Fecha" name="Fecha">
                            </div>
                        </div>

                        <h5 class="mt-4 mb-3">Moneda</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="CheckedUsd" class="form-label">Dólares (checked/"")</label>
                                <input type="text" class="form-control" id="CheckedUsd" name="CheckedUsd">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="CheckedEur" class="form-label">Euros (checked/"")</label>
                                <input type="text" class="form-control" id="CheckedEur" name="CheckedEur">
                            </div>
                        </div>

                        <h5 class="mt-4 mb-3">Información de la Transferencia</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Monto" class="form-label">Monto de la Transferencia</label>
                                <input type="text" class="form-control" id="Monto" name="Monto">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="NumeroCuenta" class="form-label">Número de Cuenta del Beneficiario</label>
                                <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="NombreBeneficiario" class="form-label">Nombre del Beneficiario</label>
                                <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="DireccionBeneficiario" class="form-label">Dirección del Beneficiario</label>
                                <input type="text" class="form-control" id="DireccionBeneficiario" name="DireccionBeneficiario">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="BeneficiaryBank" class="form-label">Nombre del Banco del Beneficiario</label>
                                <input type="text" class="form-control" id="BeneficiaryBank" name="BeneficiaryBank">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="BankAddress" class="form-label">Dirección del Banco del Beneficiario</label>
                                <input type="text" class="form-control" id="BankAddress" name="BankAddress">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Swift" class="form-label">COD.ABA - COD.SWIFT - IBAN</label>
                                <input type="text" class="form-control" id="Swift" name="Swift">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="Reference" class="form-label">Referencia</label>
                                <input type="text" class="form-control" id="Reference" name="Reference">
                            </div>
                        </div>

                        <h5 class="mt-4 mb-3">Banco Intermediario (Opcional)</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="NombreBancoIntermediario" class="form-label">Nombre del Banco Intermediario</label>
                                <input type="text" class="form-control" id="NombreBancoIntermediario" name="NombreBancoIntermediario">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="DireccionBancoIntermediario" class="form-label">Dirección del Banco Intermediario</label>
                                <input type="text" class="form-control" id="DireccionBancoIntermediario" name="DireccionBancoIntermediario">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="CodAbaCodSwiftBancoIntermediario" class="form-label">COD.ABA - COD.SWIFT del Banco Intermediario</label>
                                <input type="text" class="form-control" id="CodAbaCodSwiftBancoIntermediario" name="CodAbaCodSwiftBancoIntermediario">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="CuentaEntreBancos" class="form-label">Cuenta Entre Bancos</label>
                                <input type="text" class="form-control" id="CuentaEntreBancos" name="CuentaEntreBancos">
                            </div>
                        </div>

                        <div class="mb-3">
                            <label for="Observaciones" class="form-label">Observaciones</label>
                            <textarea class="form-control" id="Observaciones" name="Observaciones" rows="3"></textarea>
                        </div>

                        <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                            <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
            `
        break;
        case "8.html":
            text = `
            <div class="container mt-4">
    <div class="row justify-content-center">
        <div class="col-lg-8">
            <div class="card">
                <div class="card-header">
                    <h4 class="mb-0">SWIFT Receipt Form</h4>
                </div>
                <div class="card-body">
                    <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                        <div class="col-md-6 mb-3">
                                <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                        </div>
                        <div class="mb-3">
                            <label for="TextoAleatorio" class="form-label">Texto Aleatorio</label>
                            <textarea class="form-control" id="TextoAleatorio" name="TextoAleatorio" rows="10" style="font-family: monospace; white-space: pre; tab-size: 4;"></textarea>
                        </div>

                        <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                            <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
            `
        break;
        case "9.html":
            text = `
            <div class="container mt-4">
    <div class="row justify-content-center">
        <div class="col-lg-8">
            <div class="card">
                <div class="card-header">
                    <h4 class="mb-0">Column Receipt Form</h4>
                </div>
                <div class="card-body">
                    <form id="forma-libre">
                        <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Sender" class="form-label">Sender</label>
                                <input type="text" class="form-control" id="Sender" name="Sender">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="UltimoCuatroDigitos" class="form-label">Checking (últimos 4 dígitos)</label>
                                <input type="text" class="form-control" id="UltimoCuatroDigitos" name="UltimoCuatroDigitos">
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Fecha" class="form-label">Fecha</label>
                                <input type="text" class="form-control" id="Fecha" name="Fecha">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="Monto" class="form-label">Monto</label>
                                <input type="text" class="form-control" id="Monto" name="Monto">
                            </div>
                        </div>

                        <div class="mb-3">
                            <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                            <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                        </div>

                        <h5 class="mt-4 mb-3">Bank Details</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="Swift" class="form-label">SWIFT Code</label>
                                <input type="text" class="form-control" id="Swift" name="Swift">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="NumeroCuenta" class="form-label">Account Number</label>
                                <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                            </div>
                        </div>

                        <div class="mb-3">
                            <label for="Memo" class="form-label">Memo</label>
                            <input type="text" class="form-control" id="Memo" name="Memo">
                        </div>

                        <h5 class="mt-4 mb-3">Cost Breakdown</h5>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="BeneficiaryName" class="form-label">Beneficiary Name</label>
                                <input type="text" class="form-control" id="BeneficiaryName" name="BeneficiaryName">
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="TransferFees" class="form-label">Wire Cost</label>
                                <input type="text" class="form-control" id="TransferFees" name="TransferFees">
                            </div>
                        </div>

                        <div class="mb-3">
                            <label for="Total" class="form-label">Total Paid</label>
                            <input type="text" class="form-control" id="Total" name="Total">
                        </div>

                        <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                            <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
            `
        break;
        case "10.html":
            text = `
            <div class="container mt-4">
            <div class="row justify-content-center">
                <div class="col-lg-6 col-md-8">
                    <div class="card">
                        <div class="card-header">
                            <h4 class="mb-0">Payment Details Form</h4>
                        </div>
                        <div class="card-body">
                            <form id="forma-libre">
                             <input type="hidden" name="Plantilla" value="${plantilla}">
                            <input type="hidden" name="Sesion" value="${sessionHash}">
                                <div class="mb-3">
                                    <label for="MontoArriba" class="form-label">Monto</label>
                                    <input type="text" class="form-control" id="MontoArriba" name="MontoArriba">
                                </div>
    
                                <div class="mb-3">
                                    <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                    <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                                </div>
    
                                <div class="mb-3">
                                    <label for="Swift" class="form-label">SWIFT</label>
                                    <input type="text" class="form-control" id="Swift" name="Swift">
                                </div>

                                 <div class="mb-3">
                                    <label for="Imad" class="form-label">IMAD</label>
                                    <input type="text" class="form-control" id="Imad" name="Imad">
                                </div>
    
                                <div class="mb-3">
                                    <label for="NumeroCuenta" class="form-label">Número de Cuenta</label>
                                    <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                                </div>
    
                                <div class="mb-3">
                                    <label for="Fecha" class="form-label">Fecha</label>
                                    <input type="text" class="form-control" id="Fecha" name="Fecha">
                                </div>
    
                                <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                                   <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
            `
        break;
        case "11.html":
            text = `
        <div class="container mt-4">
        <div class="row justify-content-center">
            <div class="col-lg-6 col-md-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Payment Details Form</h4>
                    </div>
                    <div class="card-body">
                        <form id="forma-libre">
                         <input type="hidden" name="Plantilla" value="${plantilla}">
                        <input type="hidden" name="Sesion" value="${sessionHash}">
                            <div class="mb-3">
                                <label for="MontoArriba" class="form-label">Monto</label>
                                <input type="text" class="form-control" id="MontoArriba" name="MontoArriba">
                            </div>

                            <div class="mb-3">
                                <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                            </div>

                            <div class="mb-3">
                                <label for="RoutingNumber" class="form-label">Routing Number / ABA</label>
                                <input type="text" class="form-control" id="RoutingNumber" name="RoutingNumber">
                            </div>

                            <div class="mb-3">
                                <label for="Imad" class="form-label">IMAD</label>
                                <input type="text" class="form-control" id="Imad" name="Imad">
                            </div>

                            <div class="mb-3">
                                <label for="NumeroCuenta" class="form-label">Número de Cuenta</label>
                                <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                            </div>

                            <div class="mb-3">
                                <label for="Fecha" class="form-label">Fecha</label>
                                <input type="text" class="form-control" id="Fecha" name="Fecha">
                            </div>

                            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                               <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
        `
        break;
        default:
            text = `<div class="alert alert-danger"><b>Esta plantilla Pronto estara disponible.</b></div>`;
        break;
  }
  $("#carga").html(text);
}

function SendFormulario(e){
    let datos = formToJSON(e);
    datos = JSON.stringify(datos);
    jsonRequest(datos, "/SendFormulario", (data)=>{
        window.open("/plantilla/" + data.data, '_blank');
    })
}

function CargarFormulario(){
  let plantilla = $("#comprobante").val();
  let text = `
  <div class="form-container">
        <div class="container">
            <div class="row justify-content-center">
                <div class="col-lg-10">
                    <div class="form-card">
                        <h2 class="form-title">Remittance Receipt Form</h2>
                        <form id="forma-libre">
                        <input type="hidden" name="plantilla" value="${plantilla}">
                        <input type="hidden" name="sesion" value="${sessionHash}">
                            <h4 class="section-title">Account Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="AcountNumber" class="form-label">Account Number</label>
                                    <input type="text" class="form-control" id="AcountNumber" name="AcountNumber">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Account" class="form-label">Account</label>
                                    <input type="text" class="form-control" id="Account" name="Account">
                                </div>
                            </div>

                            <h4 class="section-title">Sender Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Sender" class="form-label">Sender</label>
                                    <input type="text" class="form-control" id="Sender" name="Sender">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="SenderName" class="form-label">Sender Name</label>
                                    <input type="text" class="form-control" id="SenderName" name="SenderName">
                                </div>
                            </div>
                            <div class="mb-3">
                                <label for="SenderAddress" class="form-label">Sender Address</label>
                                <input type="text" class="form-control" id="SenderAddress" name="SenderAddress">
                            </div>

                            <h4 class="section-title">Transfer Details</h4>
                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="TransferDate" class="form-label">Transfer Date</label>
                                    <input type="text" class="form-control" id="TransferDate" name="TransferDate">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Date" class="form-label">Date</label>
                                    <input type="text" class="form-control" id="Date" name="Date">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Time" class="form-label">Time</label>
                                    <input type="text" class="form-control" id="Time" name="Time">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Fecha" class="form-label">Fecha</label>
                                    <input type="text" class="form-control" id="Fecha" name="Fecha">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Reference" class="form-label">Reference</label>
                                    <input type="text" class="form-control" id="Reference" name="Reference">
                                </div>
                            </div>

                            <h4 class="section-title">Bank Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankName" class="form-label">Bank Name</label>
                                    <input type="text" class="form-control" id="BankName" name="BankName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Bank" class="form-label">Bank</label>
                                    <input type="text" class="form-control" id="Bank" name="Bank">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BankCountry" class="form-label">Bank Country</label>
                                    <input type="text" class="form-control" id="BankCountry" name="BankCountry">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="BankAccount" class="form-label">Bank Account</label>
                                    <input type="text" class="form-control" id="BankAccount" name="BankAccount">
                                </div>
                            </div>

                            <h4 class="section-title">Routing & Codes</h4>
                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="RoutingId" class="form-label">Routing ID</label>
                                    <input type="text" class="form-control" id="RoutingId" name="RoutingId">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="RoutingNumber" class="form-label">Routing Number</label>
                                    <input type="text" class="form-control" id="RoutingNumber" name="RoutingNumber">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Swift" class="form-label">Swift</label>
                                    <input type="text" class="form-control" id="Swift" name="Swift">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="Aba" class="form-label">ABA</label>
                                    <input type="text" class="form-control" id="Aba" name="Aba">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Imad" class="form-label">IMAD</label>
                                    <input type="text" class="form-control" id="Imad" name="Imad">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Omad" class="form-label">OMAD</label>
                                    <input type="text" class="form-control" id="Omad" name="Omad">
                                </div>
                            </div>

                            <h4 class="section-title">Amounts</h4>
                            <div class="row">
                                <div class="col-md-4 mb-3">
                                    <label for="TransferAmount" class="form-label">Transfer Amount</label>
                                    <input type="text" class="form-control" id="TransferAmount" name="TransferAmount">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="TransferFees" class="form-label">Transfer Fees</label>
                                    <input type="text" class="form-control" id="TransferFees" name="TransferFees">
                                </div>
                                <div class="col-md-4 mb-3">
                                    <label for="Total" class="form-label">Total</label>
                                    <input type="text" class="form-control" id="Total" name="Total">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Amount" class="form-label">Amount</label>
                                    <input type="text" class="form-control" id="Amount" name="Amount">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Monto" class="form-label">Monto</label>
                                    <input type="text" class="form-control" id="Monto" name="Monto">
                                </div>
                            </div>
                            <div class="mb-3">
                                <label for="MontoArriba" class="form-label">Monto Arriba</label>
                                <input type="text" class="form-control" id="MontoArriba" name="MontoArriba">
                            </div>

                            <h4 class="section-title">Beneficiary Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="AccountName" class="form-label">Account Name</label>
                                    <input type="text" class="form-control" id="AccountName" name="AccountName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="NombreBeneficiario" class="form-label">Nombre Beneficiario</label>
                                    <input type="text" class="form-control" id="NombreBeneficiario" name="NombreBeneficiario">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryName" class="form-label">Beneficiary Name</label>
                                    <input type="text" class="form-control" id="BeneficiaryName" name="BeneficiaryName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="PayeeName" class="form-label">Payee Name</label>
                                    <input type="text" class="form-control" id="PayeeName" name="PayeeName">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="Destinatario" class="form-label">Destinatario</label>
                                    <input type="text" class="form-control" id="Destinatario" name="Destinatario">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="NumeroCuenta" class="form-label">Número Cuenta</label>
                                    <input type="text" class="form-control" id="NumeroCuenta" name="NumeroCuenta">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryAccount" class="form-label">Beneficiary Account</label>
                                    <input type="text" class="form-control" id="BeneficiaryAccount" name="BeneficiaryAccount">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="DireccionBeneficiario" class="form-label">Dirección Beneficiario</label>
                                    <input type="text" class="form-control" id="DireccionBeneficiario" name="DireccionBeneficiario">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryBank" class="form-label">Beneficiary Bank</label>
                                    <input type="text" class="form-control" id="BeneficiaryBank" name="BeneficiaryBank">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="BeneficiaryAba" class="form-label">Beneficiary ABA</label>
                                    <input type="text" class="form-control" id="BeneficiaryAba" name="BeneficiaryAba">
                                </div>
                            </div>

                            <h4 class="section-title">Originator Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="OriginatorName" class="form-label">Originator Name</label>
                                    <input type="text" class="form-control" id="OriginatorName" name="OriginatorName">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="OriginatorAccount" class="form-label">Originator Account</label>
                                    <input type="text" class="form-control" id="OriginatorAccount" name="OriginatorAccount">
                                </div>
                            </div>

                            <h4 class="section-title">Intermediary Bank</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="NombreBancoIntermediario" class="form-label">Nombre Banco Intermediario</label>
                                    <input type="text" class="form-control" id="NombreBancoIntermediario" name="NombreBancoIntermediario">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="DireccionBancoIntermediario" class="form-label">Dirección Banco Intermediario</label>
                                    <input type="text" class="form-control" id="DireccionBancoIntermediario" name="DireccionBancoIntermediario">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="CodAbaCodSwiftBancoIntermediario" class="form-label">Código ABA/SWIFT Banco Intermediario</label>
                                    <input type="text" class="form-control" id="CodAbaCodSwiftBancoIntermediario" name="CodAbaCodSwiftBancoIntermediario">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="CuentaEntreBancos" class="form-label">Cuenta Entre Bancos</label>
                                    <input type="text" class="form-control" id="CuentaEntreBancos" name="CuentaEntreBancos">
                                </div>
                            </div>

                            <h4 class="section-title">Additional Information</h4>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="FedReference" class="form-label">Fed Reference</label>
                                    <input type="text" class="form-control" id="FedReference" name="FedReference">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Status" class="form-label">Status</label>
                                    <input type="text" class="form-control" id="Status" name="Status">
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6 mb-3">
                                    <label for="UltimoCuatroDigitos" class="form-label">Últimos Cuatro Dígitos</label>
                                    <input type="text" class="form-control" id="UltimoCuatroDigitos" name="UltimoCuatroDigitos">
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="Memo" class="form-label">Memo</label>
                                    <input type="text" class="form-control" id="Memo" name="Memo">
                                </div>
                            </div>

                            <h4 class="section-title">Messages</h4>
                            <div class="mb-3">
                                <label for="OriginatorToBenificiaryInformation" class="form-label">Originator to Beneficiary Information</label>
                                <input type="text" class="form-control" id="OriginatorToBenificiaryInformation" name="OriginatorToBenificiaryInformation">
                            </div>
                            <div class="mb-3">
                                <label for="MessageToPayee" class="form-label">Message to Payee</label>
                                <input type="text" class="form-control" id="MessageToPayee" name="MessageToPayee">
                            </div>
                            <div class="mb-3">
                                <label for="Observaciones" class="form-label">Observaciones</label>
                                <input type="text" class="form-control" id="Observaciones" name="Observaciones">
                            </div>
                            <div class="mb-3">
                                <label for="TextoAleatorio" class="form-label">Texto Aleatorio</label>
                                <textarea class="form-control" id="TextoAleatorio" name="TextoAleatorio" rows="4"></textarea>
                            </div>

                            <div class="text-center mt-4">
                                <button type="button" class="btn btn-primary btn-submit" onclick="SendFormulario('forma-libre')">Submit Form</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
  `
  $("#main").html(text);
}

function subirCodigo(){
  $("#main").html(`
    <div class="center" style="width: 600px;">
    <h4>Subir archivo html</h4>
    <form id="forma" method="POST" enctype="multipart/form-data">
    <input type="file" name="archivo">
    <input type="hidden" name="sesion" value="${sessionHash}">
    <br><br>
    <button class="btn btn-success" onclick="upload('forma')">Subir archivo</button>
    </form>
    </div>
    `);
}

function handleActionMenu(value) {
  switch(value) {
    case '1':
      usuarios();
      // Abrir modal de login o redirigir
      break;
    case '2':
      createUser();
      // Abrir modal de registro o redirigir
      break;
    case '3':
      location.replace('/');
      // Mostrar demo
      break;
    case '4':
      comprobantes();
      // Mostrar demo
      break;
    case 'download':
      console.log('Iniciando descarga');
      // Iniciar descarga
      break;
  }

  // Reset select
  document.querySelector('.select_menu_action').value = '';
}

function createUser(){
  $("#main").html(`
    <div class="center">
    <div class="user-form-container">
        <div class="user-form-logo">
            <!-- Puedes agregar tu logo aquí -->
            <svg width="100" height="50" viewBox="0 0 100 50" xmlns="http://www.w3.org/2000/svg">
                <rect x="10" y="10" width="80" height="30" fill="var(--color-primary)"/>
                <text x="50" y="30" text-anchor="middle" fill="white" font-size="16"></text>
            </svg>
        </div>

        <h2 class="user-form-title">Registro de Usuario</h2>

        <form id="formaRegistro">
            <div class="user-form-group">
                <label for="username" class="user-form-label">Nombre de Usuario</label>
                <input type="text" id="username" name="username" class="user-form-input" required>
            </div>

            <div class="user-form-group">
                <label for="password" class="user-form-label">Contraseña</label>
                <input type="password" id="password" name="password" class="user-form-input" required>
            </div>

            <div class="user-form-group">
                <label for="admin" class="user-form-label">Nivel de Administración</label>
                <select id="admin" name="admin" class="user-form-admin-select" required>
                    <option value="">Seleccione un nivel</option>
                    <option value="0">Usuario Normal</option>
                    <option value="1">Administrador</option>
                    <option value="2">Super Administrador</option>
                </select>
            </div>

            <div class="user-form-group">
                <label for="tlf" class="user-form-label">Teléfono</label>
                <input type="tel" id="tlf" name="tlf" class="user-form-input" required>
            </div>

            <div class="user-form-group">
                <label for="email" class="user-form-label">Correo Electrónico</label>
                <input type="email" id="email" name="email" class="user-form-input" required>
            </div>

            <input type="hidden" name="sesion" value="${sessionHash}">

            <button type="button" class="user-form-submit" onclick="accion('formaRegistro', '/create_user')">Registrar Usuario</button>
        </form>
    </div>
    </div>
    `);
}

function usuarios(){
  let datos = {
    "sesion":sessionHash
  }
  let jsonDataSerialized = JSON.stringify(datos);
  jsonRequest(jsonDataSerialized, "/users", (data)=>{
      let campos = "";
      let formas = '<form id="formaRegistro" method="POST"></form>';
      for(let i of data.data){
        let id = i.id;
        id = id.toFixed(0);
        formas += `<form id="s${id}" method="POST"></form>`
        campos+=`        
        <tr>
        <input type="hidden" form="s${id}" name="sesion" value="${sessionHash}">
        <input type="hidden" form="s${id}" name="id_user" value="${id}">
        <td><input type="text" class="form-control" name="username" form="s${id}" value="${i.username}"></td>
        <td><input type="text" class="form-control" name="tlf" form="s${id}" value="${i.telefono}"></td>
        <td>
        <select class="form-control" name="active" form="s${id}">
        <option value="${i.active}">${i.active}</option>
        <option value="0">0</option>
        <option value="1">1</option>
        </select>
        </td>
        <td>
        <select class="form-control" name="admin" form="s${id}">
        <option value="${i.admin}">${i.admin}</option>
        <option value="0">0</option>
        <option value="1">1</option>
        <option value="2">2</option>
        </select>
        </td>
        <td><input type="text" class="form-control" name="email" form="s${id}" value="${i.email}"></td>
        <td><input type="password" class="form-control" name="pass" form="s${id}" value=""></td>
        <td><button class="btn btn-success" type="button" onclick="accion2('s${id}', '/update_user')">Actualizar</button></td>
        <td><button class="btn btn-danger" type="button" onclick="accion2('s${id}', '/delete_user'), setTimeout(usuarios, 750)">Eliminar</button></td>
        </tr>
        `
      }
      let text = `
      ${formas}
      <div class="table-responsive">
      <table class="table-simple">
      <thead>
      <tr>
      <th>Username</th>
      <th>Tel</th>
      <th>Active</th>
      <th>Admin</th>
      <th>Email</th>
      <th>Password</th>
      <th>Actualizar</th>
      <th>Eliminar</th>
      </tr>
      </thead>
      <tbody>
      ${campos}
      <tr>
        <input type="hidden" form="formaRegistro" name="sesion" value="${sessionHash}">
        <td><input type="text" class="form-control" name="username" form="formaRegistro" placeholder="Nombre de usuario"></td>
        <td><input type="text" class="form-control" name="tlf" form="formaRegistro" placeholder="Telefono"></td>
        <td>
        <select class="form-control" name="active" form="formaRegistro">
        <option value="1">1</option>
        <option value="0">0</option>
        </select>
        </td>
        <td>
        <select class="form-control" name="admin" form="formaRegistro">
        <option value="0">0</option>
        <option value="1">1</option>
        <option value="2">2</option>
        </select>
        </td>
        <td><input type="text" class="form-control" name="email" form="formaRegistro" placeholder="Email"></td>
        <td><input type="password" class="form-control" name="password" form="formaRegistro" value=""></td>
        <td><button class="btn btn-success" type="button" onclick="accion2('formaRegistro', '/create_user')">Crear Nuevo</button></td>
        <td></td>
        </tr>
      </tbody>
      </table>
      </div>
      `
      $("#main").html(text);
  });
}

function menu(){
  let datos = {
      "sesion":sessionHash
  }
  let jsonDataSerialized = JSON.stringify(datos);
  jsonRequest(jsonDataSerialized, "/menu", (data)=>{
      $("#menu").html(data);
  });
}

function accion2(e, ruta){
  let jsonData = JSON.stringify(formToJSON(e));
  console.log(jsonData)
  jsonRequest(jsonData, ruta, (data, error)=>{
    if (error != null){
      return
    }
    if(data.status == "error"){
      showNotification(data.message, "danger");
      return
    }
    showNotification(data.message, "success");
  });
}

function accion(e, ruta){
  let jsonData = JSON.stringify(formToJSON(e));
  console.log(jsonData)
  jsonRequest(jsonData, ruta, (data, error)=>{
    if (error != null){
      return
    }
    if(data.status == "error"){
      showNotification(data.message, "danger");
      return
    }
    showNotification(data.message, "success");
    document.getElementById(e).reset();
  });
}

function toggleMenu() {
  const navList = document.getElementById('navList');
  navList.classList.toggle('active_menu');
}


function handleMainMenu(value) {
  // Ocultar todos los submenús
  hideAllSubmenus();

  // Mostrar submenú correspondiente
  if (value === 'services') {
    document.getElementById('servicesSubmenu').style.display = 'block';
  } else if (value === 'products') {
    document.getElementById('productsSubmenu').style.display = 'block';
  }

  // Manejar navegación
  switch(value) {
    case 'home':
      console.log('Navegando a Inicio');
      // window.location.href = '/';
      break;
    case 'about':
      console.log('Navegando a Nosotros');
      // window.location.href = '/about';
      break;
    case 'contact':
      console.log('Navegando a Contacto');
      // window.location.href = '/contact';
      break;
  }
}

function handleSubMenu(value) {
  console.log('Submenú seleccionado:', value);
  // Aquí puedes manejar la navegación de submenús
}



function hideAllSubmenus() {
  document.getElementById('servicesSubmenu').style.display = 'none';
  document.getElementById('productsSubmenu').style.display = 'none';
}