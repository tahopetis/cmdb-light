# CMDB Lite User Guide

Welcome to the CMDB Lite User Guide! This guide will help you understand how to use CMDB Lite to manage configuration items (CIs) and their relationships effectively.

## Table of Contents

- [Getting Started](#getting-started)
- [User Interface Overview](#user-interface-overview)
- [Managing Configuration Items](#managing-configuration-items)
- [Working with Relationships](#working-with-relationships)
- [Using Graph Visualization](#using-graph-visualization)
- [Managing Your Account](#managing-your-account)
- [Viewing Audit Logs](#viewing-audit-logs)
- [Common Workflows](#common-workflows)

## Getting Started

### First-Time Login

1. Open your web browser and navigate to your CMDB Lite instance (typically http://localhost:3000 for local installations).
2. You'll be presented with a login screen. Enter your username and password.
   - For first-time installations, use the default admin credentials provided during setup.
3. Click the "Login" button to access the application.

### Dashboard Overview

After logging in, you'll be taken to the Dashboard, which provides an overview of your CMDB data:

- **Total CIs**: The total number of configuration items in your CMDB.
- **CI Types**: A breakdown of CIs by type.
- **Recent Changes**: A list of recent modifications to CIs and relationships.
- **Quick Actions**: Buttons for common tasks like creating a new CI.

## User Interface Overview

The CMDB Lite interface consists of several key components:

### Main Navigation

- **Dashboard**: The home screen with an overview of your CMDB data.
- **Configuration Items**: View and manage all CIs in your system.
- **Graph**: Visualize relationships between CIs.
- **Audit Logs**: View a history of all changes made in the system.
- **Settings**: Manage your account and system preferences.

### Sidebar

The sidebar provides filtering options and quick access to different CI types:

- **All CIs**: View all configuration items regardless of type.
- **Servers**: View only server-type CIs.
- **Applications**: View only application-type CIs.
- **Databases**: View only database-type CIs.
- **Network Devices**: View only network device-type CIs.
- **Other Types**: View CIs of other custom types.

### Search Bar

The search bar at the top of the interface allows you to quickly find CIs by name or other attributes.

## Managing Configuration Items

### Creating a New CI

1. Click on the "Configuration Items" link in the main navigation.
2. Click the "Create New CI" button.
3. Fill in the CI details:
   - **Name**: A descriptive name for the CI.
   - **Type**: Select the type of CI from the dropdown (e.g., Server, Application, Database).
   - **Attributes**: Add key-value pairs to describe the CI (e.g., IP address, version, owner).
   - **Tags**: Add tags to categorize the CI (e.g., production, critical, web-tier).
4. Click "Save" to create the CI.

### Viewing and Editing CIs

1. From the CI list, click on a CI name to view its details.
2. The CI detail view shows:
   - Basic information (name, type, creation date)
   - Attributes and tags
   - Relationships to other CIs
   - Recent changes
3. To edit a CI, click the "Edit" button and modify the desired fields.
4. Click "Save" to apply your changes.

### Deleting CIs

1. From the CI list or detail view, click the "Delete" button.
2. Confirm the deletion in the popup dialog.
3. Note: Deleting a CI will also remove all relationships associated with it.

### Searching and Filtering CIs

1. Use the search bar to find CIs by name or attribute.
2. Use the sidebar to filter CIs by type.
3. Combine search and filters for more specific results.

## Working with Relationships

### Creating Relationships

1. Navigate to the detail view of the source CI.
2. Click the "Add Relationship" button.
3. Select the target CI from the dropdown list.
4. Choose the relationship type (e.g., "depends on", "hosts", "connects to").
5. Click "Save" to create the relationship.

### Viewing Relationships

1. From the CI detail view, scroll to the "Relationships" section.
2. This section shows:
   - **Outgoing Relationships**: Relationships from this CI to other CIs.
   - **Incoming Relationships**: Relationships from other CIs to this CI.
3. Click on any related CI to view its details.

### Managing Relationships

1. From the CI detail view, find the relationship you want to modify.
2. Click the "Edit" button next to the relationship to change its type.
3. Click the "Delete" button to remove the relationship.
4. Confirm your action in the popup dialog.

## Using Graph Visualization

### Accessing the Graph View

1. Click on the "Graph" link in the main navigation.
2. The graph view shows all CIs and their relationships as a network diagram.

### Navigating the Graph

- **Pan**: Click and drag to move around the graph.
- **Zoom**: Use the mouse wheel or pinch gesture to zoom in and out.
- **Select CI**: Click on a CI node to see its details in a side panel.
- **Expand/Collapse**: Double-click on a CI to expand or collapse its relationships.

### Customizing the Graph View

1. Use the toolbar at the top of the graph view to:
   - Filter by CI type
   - Change the layout algorithm
   - Adjust the zoom level
   - Reset the view
2. Use the search bar to find and highlight specific CIs in the graph.

## Managing Your Account

### Changing Your Password

1. Click on your username in the top-right corner and select "Settings".
2. Navigate to the "Account" tab.
3. Enter your current password and new password.
4. Click "Update Password" to save your changes.

### Updating Your Profile

1. Click on your username in the top-right corner and select "Settings".
2. Navigate to the "Profile" tab.
3. Update your information as needed.
4. Click "Save Changes" to apply your updates.

## Viewing Audit Logs

### Accessing Audit Logs

1. Click on the "Audit Logs" link in the main navigation.
2. The audit log view shows a history of all changes made in the system.

### Filtering Audit Logs

1. Use the filter options to narrow down the log entries:
   - **Date Range**: Select a specific time period.
   - **User**: Filter by the user who made the changes.
   - **Action Type**: Filter by type of action (create, update, delete).
   - **Entity Type**: Filter by type of entity (CI, relationship).
2. Click "Apply Filters" to update the view.

### Understanding Audit Log Entries

Each audit log entry includes:
- **Timestamp**: When the action occurred.
- **User**: Who performed the action.
- **Action**: What was done (create, update, delete).
- **Entity Type**: The type of entity affected (CI, relationship).
- **Entity Name**: The name of the affected entity.
- **Details**: Additional information about the changes made.

## Common Workflows

### Tracking Service Dependencies

1. Create CIs for all components of your service (servers, applications, databases).
2. Create relationships between these CIs to represent dependencies.
3. Use the graph view to visualize the entire service architecture.
4. Regularly review and update the relationships as your service evolves.

### Impact Analysis

1. When planning a change to a CI, first view its relationships.
2. Use the graph view to identify all dependent CIs.
3. Consider the impact of your change on all dependent components.
4. Update relevant documentation and notify affected teams.

### Onboarding New Team Members

1. Use the graph view to give new team members an overview of your infrastructure.
2. Show them how to find and update CI information.
3. Encourage them to keep the CMDB data accurate as they make changes.

### Troubleshooting Issues

1. When an issue occurs, use the CMDB to understand the affected components.
2. Check the relationships to identify potential root causes.
3. Review the audit logs to see recent changes that might be related to the issue.
4. Use this information to guide your troubleshooting process.

## Tips for Effective CMDB Usage

- **Keep it Current**: Make updating the CMDB part of your change management process.
- **Be Consistent**: Use consistent naming conventions and attribute values.
- **Document Relationships**: Ensure all important dependencies are captured.
- **Review Regularly**: Schedule regular reviews to clean up outdated information.
- **Integrate with Workflows**: Use the CMDB data to inform other processes like incident management and capacity planning.

## Getting Help

If you encounter issues or have questions about using CMDB Lite:
- Check the [FAQ section](faq.md) for answers to common questions.
- Contact your system administrator for account or access issues.
- For feature requests or bug reports, see the [Contribution Guide](../project/contributing.md).