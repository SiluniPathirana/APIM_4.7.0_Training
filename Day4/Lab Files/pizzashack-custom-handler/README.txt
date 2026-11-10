PizzaShack Custom Handler (Lab 19)

Build
  mvn clean package
  Output: target/pizzashack-custom-handler-1.0.jar

Deploy
  1. Copy target/pizzashack-custom-handler-1.0.jar to <APIM_HOME>/repository/components/lib
  2. Edit <APIM_HOME>/repository/resources/api_templates/velocity_template.xml (see Lab 19, steps 6-8)
  3. Restart API Manager and deploy a new revision of the PizzaShack API

Requirements: JDK 21 and Maven. The first build downloads synapse-core from the WSO2 public Maven repository.
